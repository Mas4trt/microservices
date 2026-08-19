package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	inventoryAPI "github.com/Mas4trt/microservices/inventory/internal/api/inventory/v1"
	"github.com/Mas4trt/microservices/inventory/internal/config"
	"github.com/Mas4trt/microservices/inventory/internal/model"
	inventoryRepository "github.com/Mas4trt/microservices/inventory/internal/repository/part"
	inventoryService "github.com/Mas4trt/microservices/inventory/internal/service/part"
	inventoryV1 "github.com/Mas4trt/microservices/shared/pkg/proto/inventory/v1"
)

const configPath = "./deploy/compose/inventory/.env"

func seedParts() map[string]model.Part {
	now := time.Now()

	return map[string]model.Part{
		"550e8400-e29b-41d4-a716-446655440000": {
			UUID:          "550e8400-e29b-41d4-a716-446655440000",
			Name:          "Main Booster Engine",
			Description:   "Основной маршевый двигатель первой ступени",
			Price:         1_250_000.0,
			StockQuantity: 4,
			Category:      model.CategoryEngine,

			Dimensions: model.Dimensions{
				Length: 320,
				Width:  180,
				Height: 180,
				Weight: 2100,
			},

			Manufacturer: model.Manufacturer{
				Name:    "RocketDyne Corp",
				Country: "Germany",
				Website: "https://rocketdyne.example.com",
			},

			Tags: []string{
				"main",
				"booster",
			},

			Metadata: map[string]any{
				"thrust_kn": 7_600.5,
				"reusable":  true,
			},

			CreatedAt: now,
			UpdatedAt: &now,
		},

		"f47ac10b-58cc-4372-a567-0e02b2c3d479": {
			UUID:          "f47ac10b-58cc-4372-a567-0e02b2c3d479",
			Name:          "Cryo Fuel Tank",
			Description:   "Криогенный топливный бак",
			Price:         340_000.0,
			StockQuantity: 12,
			Category:      model.CategoryFuel,

			Dimensions: model.Dimensions{
				Length: 500,
				Width:  220,
				Height: 220,
				Weight: 900,
			},

			Manufacturer: model.Manufacturer{
				Name:    "TankWorks Inc",
				Country: "USA",
				Website: "https://tankworks.example.com",
			},

			Tags: []string{
				"secondary",
				"cryogenic",
			},

			Metadata: map[string]any{
				"capacity_l": int64(18_000),
			},

			CreatedAt: now,
			UpdatedAt: &now,
		},
	}
}

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
	if err != nil {
		log.Printf("failed to connect to MongoDB: %v\n", err)
		return
	}
	defer func() {
		cerr := client.Disconnect(context.Background())
		if cerr != nil {
			log.Printf("failed to disconnect from MongoDB: %v\n", cerr)
		}
	}()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("failed to ping MongoDB: %v\n", err)
		return
	}
	log.Println("✅ Connected to MongoDB")

	lis, err := net.Listen("tcp", config.AppConfig().InternalGRPC.Address())
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	s := grpc.NewServer()

	repo := inventoryRepository.NewRepository(client, config.AppConfig().Mongo.DatabaseName())
	service := inventoryService.NewService(repo)
	api := inventoryAPI.NewAPI(service)

	if err := repo.Init(ctx, seedParts()); err != nil {
		log.Printf("Fatal error: %v", err)
	}

	inventoryV1.RegisterInventoryServiceServer(s, api)

	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on %s\n", config.AppConfig().InternalGRPC.Address())
		if err := s.Serve(lis); err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("gRPC server stopped.")
}
