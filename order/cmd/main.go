package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ogen-go/ogen/validate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/wrapperspb"

	orderV1 "github.com/Mas4trt/microservices/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/Mas4trt/microservices/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/Mas4trt/microservices/shared/pkg/proto/payment/v1"
)

const (
	httpPort          = "8080"
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]*orderV1.OrderDto
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[uuid.UUID]*orderV1.OrderDto),
	}
}

func (s *OrderStorage) CreateOrder(uuid uuid.UUID, order *orderV1.OrderDto) error {
	if order == nil {
		return errors.New("order cannot be nil")
	}

	orderCopy := cloneOrder(order)

	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[uuid] = orderCopy

	return nil
}

func (s *OrderStorage) GetOrder(uuid uuid.UUID) *orderV1.OrderDto {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[uuid]
	if !ok {
		return nil
	}

	return cloneOrder(order)
}

func (s *OrderStorage) UpdateOrder(uuid uuid.UUID, order *orderV1.OrderDto) error {
	if order == nil {
		return errors.New("order cannot be nil")
	}

	orderCopy := cloneOrder(order)

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.orders[uuid]; !exists {
		return fmt.Errorf("order %s not found", uuid)
	}

	s.orders[uuid] = orderCopy

	return nil
}

func cloneOrder(order *orderV1.OrderDto) *orderV1.OrderDto {
	if order == nil {
		return nil
	}

	clone := *order

	if order.PartUuids != nil {
		clone.PartUuids = append([]uuid.UUID(nil), order.PartUuids...)
	}

	return &clone
}

type OrderHandler struct {
	inventoryClient inventoryV1.InventoryServiceClient
	paymentClient   paymentV1.PaymentServiceClient
	storage         *OrderStorage
}

func NewOrderHandler(storage *OrderStorage, inventoryClient inventoryV1.InventoryServiceClient, paymentClient paymentV1.PaymentServiceClient) *OrderHandler {
	return &OrderHandler{
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
		storage:         storage,
	}
}

func (h *OrderHandler) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	order := h.storage.GetOrder(params.OrderUUID)
	if order == nil {
		return &orderV1.CancelOrderNotFound{
			Code:    "ORDER_NOT_FOUND",
			Message: "order not found. orderUUID: " + params.OrderUUID.String(),
		}, nil
	}

	switch order.Status {
	case orderV1.OrderStatusPENDINGPAYMENT:
		order.Status = orderV1.OrderStatusCANCELLED

		if err := h.storage.UpdateOrder(params.OrderUUID, order); err != nil {
			return &orderV1.CancelOrderInternalServerError{
				Code:    "INTERNAL_ERROR",
				Message: "failed to update order",
			}, nil
		}

		return &orderV1.CancelOrderNoContent{}, nil
	case orderV1.OrderStatusPAID:
		return &orderV1.CancelOrderConflict{
			Code:    "ORDER_ALREADY_PAID",
			Message: "order is already paid and cannot be cancelled. orderUUID: " + params.OrderUUID.String(),
		}, nil
	case orderV1.OrderStatusCANCELLED:
		return &orderV1.CancelOrderConflict{
			Code:    "ORDER_ALREADY_CANCELLED",
			Message: "order is already cancelled. orderUUID: " + params.OrderUUID.String(),
		}, nil
	default:
		return &orderV1.CancelOrderConflict{
			Code:    "INVALID_ORDER_STATUS",
			Message: "order cannot be cancelled in current status",
		}, nil
	}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	// If req = nil Validate return validate.ErrNilPointer
	if err := req.Validate(); err != nil {
		return validationError(err), nil
	}

	uuids := make([]*wrapperspb.StringValue, 0, len(req.PartUuids))

	for _, uuid := range req.PartUuids {
		uuids = append(uuids, wrapperspb.String(uuid.String()))
	}

	resp, err := h.inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{
			Uuids: uuids,
		},
	})
	if err != nil {
		return &orderV1.CreateOrderBadGateway{
			Code:    "UPSTREAM_ERROR",
			Message: "inventory service is unavailable",
		}, nil
	}

	parts := resp.GetParts()

	foundParts := make(map[string]struct{}, len(parts))

	for _, part := range parts {
		foundParts[part.GetUuid()] = struct{}{}
	}

	for _, requestedUUID := range req.GetPartUuids() {
		if _, ok := foundParts[requestedUUID.String()]; !ok {
			return &orderV1.CreateOrderNotFound{
				Code:    "PART_NOT_FOUND",
				Message: "part not found. partUUID: " + requestedUUID.String(),
			}, nil
		}
	}

	var totalPrice float64

	for _, part := range parts {
		totalPrice += part.Price
	}

	orderUUID, err := uuid.NewRandom()
	if err != nil {
		return &orderV1.CreateOrderInternalServerError{
			Code:    "INTERNAL_ERROR",
			Message: "failed to generate order UUID",
		}, nil
	}

	newOrder := &orderV1.OrderDto{
		OrderUUID:  orderUUID,
		UserUUID:   req.UserUUID,
		PartUuids:  req.PartUuids,
		TotalPrice: totalPrice,
		Status:     orderV1.OrderStatusPENDINGPAYMENT,
	}

	err = h.storage.CreateOrder(orderUUID, newOrder)
	if err != nil {
		return &orderV1.CreateOrderInternalServerError{
			Code:    "INTERNAL_ERROR",
			Message: "failed to save order",
		}, nil
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: totalPrice,
	}, nil
}

func (h *OrderHandler) GetOrder(ctx context.Context, params orderV1.GetOrderParams) (orderV1.GetOrderRes, error) {
	order := h.storage.GetOrder(params.OrderUUID)
	if order == nil {
		return &orderV1.GetOrderNotFound{
			Code:    "ORDER_NOT_FOUND",
			Message: "order not found. orderUUID: " + params.OrderUUID.String(),
		}, nil
	}

	return order, nil
}

func (h *OrderHandler) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	// If req = nil Validate return validate.ErrNilPointer
	if err := req.Validate(); err != nil {
		return validationError(err), nil
	}

	order := h.storage.GetOrder(params.OrderUUID)
	if order == nil {
		return &orderV1.PayOrderNotFound{
			Code:    "ORDER_NOT_FOUND",
			Message: "order not found. orderUUID: " + params.OrderUUID.String(),
		}, nil
	}

	switch order.Status {
	case orderV1.OrderStatusPAID:
		return &orderV1.PayOrderConflict{
			Code:    "ORDER_ALREADY_PAID",
			Message: "order is already paid",
		}, nil

	case orderV1.OrderStatusCANCELLED:
		return &orderV1.PayOrderConflict{
			Code:    "ORDER_CANCELLED",
			Message: "order is cancelled and cannot be paid",
		}, nil

	case orderV1.OrderStatusPENDINGPAYMENT:
		// Order can be paid.

	default:
		return &orderV1.PayOrderConflict{
			Code:    "INVALID_ORDER_STATUS",
			Message: "order cannot be paid in current status",
		}, nil
	}

	paymentMethod, err := PaymentMethodToProto(req.PaymentMethod)
	if err != nil {
		return &orderV1.ValidationError{
			Code:    "VALIDATION_ERROR",
			Message: "invalid payment method",
			Violations: []orderV1.ValidationErrorViolationsItem{
				{
					Field:   "payment_method",
					Message: err.Error(),
				},
			},
		}, nil
	}

	resp, err := h.paymentClient.PayOrder(ctx, &paymentV1.PayOrderRequest{
		UserUuid:      order.UserUUID.String(),
		OrderUuid:     order.OrderUUID.String(),
		PaymentMethod: paymentMethod,
	})
	if err != nil {
		return &orderV1.PayOrderBadGateway{
			Code:    "UPSTREAM_ERROR",
			Message: "payment service is unavailable",
		}, nil
	}

	transactionUUID, err := uuid.Parse(resp.GetTransactionUuid())
	if err != nil {
		return &orderV1.PayOrderInternalServerError{
			Code:    "INTERNAL_ERROR",
			Message: "payment service returned invalid transaction UUID",
		}, nil
	}

	order.TransactionUUID = orderV1.OptNilUUID{
		Value: transactionUUID,
	}

	order.PaymentMethod = orderV1.NewOptNilPaymentMethod(req.PaymentMethod)

	order.Status = orderV1.OrderStatusPAID

	err = h.storage.UpdateOrder(params.OrderUUID, order)
	if err != nil {
		return &orderV1.PayOrderInternalServerError{
			Code:    "INTERNAL_ERROR",
			Message: "failed to update order",
		}, nil
	}

	return &orderV1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil
}

func (h *OrderHandler) NewError(ctx context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: 500,
		Response: orderV1.GenericError{
			Code:    "INTERNAL_ERROR",
			Message: "internal server error",
		},
	}
}

func PaymentMethodToProto(
	method orderV1.PaymentMethod,
) (paymentV1.PaymentMethod, error) {
	switch method {
	case orderV1.PaymentMethodCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CARD, nil

	case orderV1.PaymentMethodSBP:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_SBP, nil

	case orderV1.PaymentMethodCREDITCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD, nil

	case orderV1.PaymentMethodINVESTORMONEY:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY, nil

	default:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED,
			fmt.Errorf("unsupported payment method: %q", method)
	}
}

func validationError(err error) *orderV1.ValidationError {
	var validationErr *validate.Error

	if errors.As(err, &validationErr) {
		violations := make(
			[]orderV1.ValidationErrorViolationsItem,
			0,
			len(validationErr.Fields),
		)

		for _, field := range validationErr.Fields {
			violations = append(
				violations,
				orderV1.ValidationErrorViolationsItem{
					Field:   field.Name,
					Message: field.Error.Error(),
				},
			)
		}

		return &orderV1.ValidationError{
			Code:       "VALIDATION_ERROR",
			Message:    "request validation failed",
			Violations: violations,
		}
	}

	return &orderV1.ValidationError{
		Code:    "VALIDATION_ERROR",
		Message: err.Error(),
	}
}

func main() {
	inventoryConn, err := grpc.NewClient(
		"localhost:50050",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer inventoryConn.Close()

	inventoryClient := inventoryV1.NewInventoryServiceClient(inventoryConn)

	paymentConn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer paymentConn.Close()

	paymentClient := paymentV1.NewPaymentServiceClient(paymentConn)

	storage := NewOrderStorage()
	orderHandler := NewOrderHandler(storage, inventoryClient, paymentClient)

	orderServer, err := orderV1.NewServer(orderHandler)
	if err != nil {
		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
