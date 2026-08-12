package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	inventoryV1 "github.com/Mas4trt/microservices/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50050

type InventoryService struct {
	inventoryV1.UnimplementedInventoryServiceServer

	mu    sync.RWMutex
	parts map[string]*inventoryV1.Part
}

func (s *InventoryService) GetPart(_ context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	part, ok := s.parts[req.GetUuid()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", req.GetUuid())
	}

	return &inventoryV1.GetPartResponse{
		Part: part,
	}, nil
}

func (s *InventoryService) ListParts(_ context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	filter := req.GetFilter()

	uuids := stringValueSet(filter.GetUuids())
	names := stringValueSet(filter.GetNames())
	countries := stringValueSet(filter.GetManufacturerCountries())
	tags := stringValueSet(filter.GetTags())
	categories := categorySet(filter.GetCategories())

	result := make([]*inventoryV1.Part, 0, len(s.parts))

	for _, part := range s.parts {
		if !matchesSet(uuids, part.GetUuid()) {
			continue
		}
		if !matchesSet(names, part.GetName()) {
			continue
		}
		if !matchesCategorySet(categories, part.GetCategory()) {
			continue
		}
		if !matchesSet(countries, part.GetManufacturer().GetCountry()) {
			continue
		}
		if !matchesAnySet(tags, part.GetTags()) {
			continue
		}

		result = append(result, part)
	}

	return &inventoryV1.ListPartsResponse{
		Parts: result,
	}, nil
}

// stringValueSet строит множество значений из []*wrapperspb.StringValue.
// Пустой/nil-список -> пустое множество -> фильтр по этому полю не применяется.
func stringValueSet(values []*wrapperspb.StringValue) map[string]struct{} {
	if len(values) == 0 {
		return nil
	}

	set := make(map[string]struct{}, len(values))
	for _, v := range values {
		set[v.GetValue()] = struct{}{}
	}

	return set
}

func categorySet(values []inventoryV1.Category) map[inventoryV1.Category]struct{} {
	if len(values) == 0 {
		return nil
	}

	set := make(map[inventoryV1.Category]struct{}, len(values))
	for _, v := range values {
		set[v] = struct{}{}
	}

	return set
}

// matchesSet: nil/пустое множество означает "фильтр не задан" -> всегда true.
func matchesSet(set map[string]struct{}, value string) bool {
	if set == nil {
		return true
	}

	_, ok := set[value]
	return ok
}

func matchesCategorySet(set map[inventoryV1.Category]struct{}, value inventoryV1.Category) bool {
	if set == nil {
		return true
	}

	_, ok := set[value]
	return ok
}

// matchesAnySet: ИЛИ внутри поля — хотя бы один тег детали входит в фильтр.
func matchesAnySet(set map[string]struct{}, values []string) bool {
	if set == nil {
		return true
	}

	for _, v := range values {
		if _, ok := set[v]; ok {
			return true
		}
	}

	return false
}

func seedParts() map[string]*inventoryV1.Part {
	now := timestamppb.Now()

	return map[string]*inventoryV1.Part{
		"550e8400-e29b-41d4-a716-446655440000": {
			Uuid:          "550e8400-e29b-41d4-a716-446655440000",
			Name:          "Main Booster Engine",
			Description:   "Основной маршевый двигатель первой ступени",
			Price:         1_250_000.0,
			StockQuantity: 4,
			Category:      inventoryV1.Category_CATEGORY_ENGINE,
			Dimensions: &inventoryV1.Dimensions{
				Length: 320, Width: 180, Height: 180, Weight: 2100,
			},
			Manufacturer: &inventoryV1.Manufacturer{
				Name:    "RocketDyne Corp",
				Country: "Germany",
				Website: "https://rocketdyne.example.com",
			},
			Tags: []string{"main", "booster"},
			Metadata: map[string]*inventoryV1.Value{
				"thrust_kn": {Value: &inventoryV1.Value_DoubleValue{DoubleValue: 7600.5}},
				"reusable":  {Value: &inventoryV1.Value_BoolValue{BoolValue: true}},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		"f47ac10b-58cc-4372-a567-0e02b2c3d479": {
			Uuid:          "f47ac10b-58cc-4372-a567-0e02b2c3d479",
			Name:          "Cryo Fuel Tank",
			Description:   "Криогенный топливный бак",
			Price:         340_000.0,
			StockQuantity: 12,
			Category:      inventoryV1.Category_CATEGORY_FUEL,
			Dimensions: &inventoryV1.Dimensions{
				Length: 500, Width: 220, Height: 220, Weight: 900,
			},
			Manufacturer: &inventoryV1.Manufacturer{
				Name:    "TankWorks Inc",
				Country: "USA",
				Website: "https://tankworks.example.com",
			},
			Tags: []string{"secondary", "cryogenic"},
			Metadata: map[string]*inventoryV1.Value{
				"capacity_l": {Value: &inventoryV1.Value_Int64Value{Int64Value: 18000}},
			},
			CreatedAt: now,
			UpdatedAt: now,
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

	service := &InventoryService{
		parts: seedParts(),
	}

	inventoryV1.RegisterInventoryServiceServer(s, service)

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
