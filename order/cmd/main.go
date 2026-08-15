package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderHandler "github.com/Mas4trt/microservices/order/internal/api/order/v1"
	inventoryClient "github.com/Mas4trt/microservices/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/Mas4trt/microservices/order/internal/client/grpc/payment/v1"
	orderRepository "github.com/Mas4trt/microservices/order/internal/repository/order"
	orderService "github.com/Mas4trt/microservices/order/internal/service/order"
	orderV1 "github.com/Mas4trt/microservices/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/Mas4trt/microservices/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/Mas4trt/microservices/shared/pkg/proto/payment/v1"
)

const (
	httpPort          = "8081"
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func main() {
	inventoryConn, err := grpc.NewClient(
		"localhost:50050",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("Error: %v", err)
	}
	defer func() {
		if err := inventoryConn.Close(); err != nil {
			log.Printf("failed to close inventory connection: %v", err)
		}
	}()

	inventoryGRPCClient := inventoryV1.NewInventoryServiceClient(inventoryConn)
	invClient := inventoryClient.NewClient(inventoryGRPCClient)

	paymentConn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("Error: %v", err)
	}
	defer func() {
		if err := paymentConn.Close(); err != nil {
			log.Printf("failed to close payment connection: %v", err)
		}
	}()

	paymentGRPCClient := paymentV1.NewPaymentServiceClient(paymentConn)
	payClient := paymentClient.NewClient(paymentGRPCClient)

	repo := orderRepository.NewRepository()
	orderSvc := orderService.NewService(repo, invClient, payClient)
	orderHandler := orderHandler.NewOrderHandler(orderSvc)

	orderServer, err := orderV1.NewServer(orderHandler)
	if err != nil {
		// fmt.Printf("ошибка создания сервера OpenAPI")
		log.Printf("ошибка создания сервера OpenAPI")
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
