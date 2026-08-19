package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderHandler "github.com/Mas4trt/microservices/order/internal/api/order/v1"
	inventoryClient "github.com/Mas4trt/microservices/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/Mas4trt/microservices/order/internal/client/grpc/payment/v1"
	"github.com/Mas4trt/microservices/order/internal/config"
	"github.com/Mas4trt/microservices/order/internal/migrator"
	orderRepository "github.com/Mas4trt/microservices/order/internal/repository/order"
	orderService "github.com/Mas4trt/microservices/order/internal/service/order"
	orderV1 "github.com/Mas4trt/microservices/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/Mas4trt/microservices/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/Mas4trt/microservices/shared/pkg/proto/payment/v1"
)

const (
	configPath      = "./deploy/compose/order/.env"
	shutdownTimeout = 10 * time.Second
)

func main() {
	if err := config.Load(configPath); err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	pgxCfg, err := pgxpool.ParseConfig(config.AppConfig().Postgres.URI())
	if err != nil {
		log.Printf("failed to parse DB_URI: %v\n", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

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

	migratorRunner := migrator.NewMigrator(migrationDB, config.AppConfig().Postgres.MigrationDir())
	if err := migratorRunner.Up(); err != nil {
		log.Printf("Ошибка миграции базы данных: %v\n", err)
		return
	}
	log.Println("✅ Migrations applied")

	inventoryConn, err := grpc.NewClient(
		config.AppConfig().InventoryGRPC.Address(),
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
		config.AppConfig().PaymentGRPC.Address(),
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
		Addr:              config.AppConfig().OrderHTTP.Address(),
		Handler:           r,
		ReadHeaderTimeout: config.AppConfig().OrderHTTP.ReadTimeout(),
	}

	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", config.AppConfig().OrderHTTP.Address())
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctxShutdown)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
