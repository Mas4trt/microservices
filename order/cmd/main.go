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
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderHandler "github.com/Mas4trt/microservices/order/internal/api/order/v1"
	inventoryClient "github.com/Mas4trt/microservices/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/Mas4trt/microservices/order/internal/client/grpc/payment/v1"
	"github.com/Mas4trt/microservices/order/internal/migrator"
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
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("failed to load .env file: %v\n", err)
	}

	dbURI := os.Getenv("DB_URI")
	migrationsDir := os.Getenv("MIGRATIONS_DIR")

	pgxCfg, err := pgxpool.ParseConfig(dbURI)
	if err != nil {
		log.Printf("failed to parse DB_URI: %v\n", err)
		return
	}

	ctx := context.Background()

	pool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		log.Printf("failed to connect to database: %v\n", err)
		return
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Printf("failed to ping postgres: %v", err)
		return
	}
	log.Println("✅ Connected to Postgres")

	migrationDB := stdlib.OpenDB(*pgxCfg.ConnConfig.Copy())
	defer func() {
		if err := migrationDB.Close(); err != nil {
			log.Printf("failed to close migration db connection: %v\n", err)
		}
	}()

	migratorRunner := migrator.NewMigrator(migrationDB, migrationsDir)
	if err := migratorRunner.Up(); err != nil {
		log.Printf("Ошибка миграции базы данных: %v\n", err)
		return
	}
	log.Println("✅ Migrations applied")

	inventoryConn, err := grpc.NewClient(
		"localhost:50050",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to create inventory gRPC client: %v", err)
		return
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
		log.Printf("failed to create payment gRPC client: %v", err)
		return
	}
	defer func() {
		if err := paymentConn.Close(); err != nil {
			log.Printf("failed to close payment connection: %v", err)
		}
	}()

	paymentGRPCClient := paymentV1.NewPaymentServiceClient(paymentConn)
	payClient := paymentClient.NewClient(paymentGRPCClient)

	repo := orderRepository.NewRepository(pool)
	orderSvc := orderService.NewService(repo, invClient, payClient)
	orderHandler := orderHandler.NewOrderHandler(orderSvc)

	orderServer, err := orderV1.NewServer(orderHandler)
	if err != nil {
		log.Printf("ошибка создания сервера OpenAPI: %v", err)
		return
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
		err := server.ListenAndServe()
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
