package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"

	paymentV1 "github.com/Mas4trt/microservices/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 50051

type paymentService struct {
	paymentV1.UnimplementedPaymentServiceServer
}

func (s *paymentService) PayOrder(_ context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	transactionUUID := uuid.NewString()

	log.Printf(
		"Processed payment for order %s, transaction ID: %s",
		req.GetOrderUuid(),
		transactionUUID,
	)

	return &paymentV1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil
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

	service := &paymentService{}

	paymentV1.RegisterPaymentServiceServer(s, service)

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
