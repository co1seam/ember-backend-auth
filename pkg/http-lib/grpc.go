package http_lib

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// GRPCServer представляет gRPC сервер
type GRPCServer struct {
	// Встроенный gRPC сервер
	server *grpc.Server
	
	// Конфигурация сервера
	config *ServerConfig
	
	// Список зарегистрированных сервисов
	services []interface{}
}

// NewGRPCServer создает новый экземпляр gRPC сервера
func NewGRPCServer(config *ServerConfig) *GRPCServer {
	// Создаем опции для gRPC сервера
	opts := []grpc.ServerOption{
		grpc.MaxConcurrentStreams(config.MaxConcurrentStreams),
		grpc.MaxRecvMsgSize(config.MaxRecvMsgSize),
		grpc.MaxSendMsgSize(config.MaxSendMsgSize),
		grpc.InitialWindowSize(config.InitialWindowSize),
		grpc.InitialConnWindowSize(config.InitialConnWindowSize),
	}

	// Добавляем middleware, если они есть
	if len(config.GRPCMiddleware) > 0 {
		opts = append(opts, grpc.ChainUnaryInterceptor(config.GRPCMiddleware...))
	}

	// Создаем gRPC сервер
	server := grpc.NewServer(opts...)

	return &GRPCServer{
		server:   server,
		config:   config,
		services: make([]interface{}, 0),
	}
}

// RegisterService регистрирует gRPC сервис
func (s *GRPCServer) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {
	s.server.RegisterService(desc, impl)
	s.services = append(s.services, impl)
}

// Start запускает gRPC сервер
func (s *GRPCServer) Start() error {
	// Создаем TCP слушатель
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.config.Host, s.config.Port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	// Включаем reflection для отладки
	if s.config.EnableReflection {
		reflection.Register(s.server)
	}

	// Запускаем сервер
	if err := s.server.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

// Stop останавливает gRPC сервер
func (s *GRPCServer) Stop() {
	if s.server != nil {
		s.server.GracefulStop()
	}
}

// GetServer возвращает внутренний gRPC сервер
func (s *GRPCServer) GetServer() *grpc.Server {
	return s.server
}

// GetServices возвращает список зарегистрированных сервисов
func (s *GRPCServer) GetServices() []interface{} {
	return s.services
} 