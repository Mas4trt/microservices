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
		log.Fatal(err)
	}
	defer inventoryConn.Close()

	inventoryGRPCClient := inventoryV1.NewInventoryServiceClient(inventoryConn)
	inventoryClient := inventoryClient.NewClient(inventoryGRPCClient)

	paymentConn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer paymentConn.Close()

	paymentGRPCClient := paymentV1.NewPaymentServiceClient(paymentConn)
	paymentClient := paymentClient.NewClient(paymentGRPCClient)

	repo := orderRepository.NewRepository()
	service := orderService.NewService(repo, inventoryClient, paymentClient)
	orderHandler := orderHandler.NewOrderHandler(service)

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
