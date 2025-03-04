package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

// server реализует gRPC сервис
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello реализует метод gRPC сервиса
func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: fmt.Sprintf("Hello %s!", req.GetName())}, nil
}

func main() {
	// Создаем конфигурацию сервера
	config := http_lib.DefaultServerConfig()
	config.Port = 50051 // Стандартный порт для gRPC

	// Создаем gRPC сервер
	grpcServer := http_lib.NewGRPCServer(config)

	// Регистрируем сервис
	pb.RegisterGreeterServer(grpcServer.GetServer(), &server{})

	// Включаем reflection для отладки
	reflection.Register(grpcServer.GetServer())

	// Запускаем сервер
	log.Printf("Starting gRPC server on %s:%d", config.Host, config.Port)
	if err := grpcServer.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
} 