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

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	inventoryAPI "github.com/Mas4trt/microservices/inventory/internal/api/inventory/v1"
	"github.com/Mas4trt/microservices/inventory/internal/model"
	inventoryRepository "github.com/Mas4trt/microservices/inventory/internal/repository/part"
	inventoryService "github.com/Mas4trt/microservices/inventory/internal/service/part"
	inventoryV1 "github.com/Mas4trt/microservices/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50050

func seedParts() map[string]model.Part {
	now := time.Now()

	return map[string]model.Part{
		"550e8400-e29b-41d4-a716-446655440000": {
			Uuid:          "550e8400-e29b-41d4-a716-446655440000",
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
			Uuid:          "f47ac10b-58cc-4372-a567-0e02b2c3d479",
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
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
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

	repo := inventoryRepository.NewRepository()
	service := inventoryService.NewService(repo)
	api := inventoryAPI.NewAPI(service)

	ctx := context.Background()
	if err := repo.Init(ctx, seedParts()); err != nil {
		log.Printf("Fatal error: %v", err)
	}

	inventoryV1.RegisterInventoryServiceServer(s, api)

	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
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
