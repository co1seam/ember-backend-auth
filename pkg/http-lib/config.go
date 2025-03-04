package http_lib

import (
	"time"
	"google.golang.org/grpc"
)

// Config содержит настройки HTTP сервера
type Config struct {
	Port           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
}

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() Config {
	return Config{
		Port:           "8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1024,
	}
}

// ServerConfig представляет конфигурацию сервера
type ServerConfig struct {
	// HTTP настройки
	Host         string
	Port         int
	ReadTimeout  int
	WriteTimeout int
	IdleTimeout  int
	MaxHeaderBytes int
	
	// gRPC настройки
	MaxConcurrentStreams    uint32
	MaxRecvMsgSize         int
	MaxSendMsgSize         int
	InitialWindowSize      int32
	InitialConnWindowSize  int32
	EnableReflection       bool
	GRPCMiddleware         []grpc.UnaryServerInterceptor
}

// DefaultServerConfig возвращает конфигурацию сервера по умолчанию
func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		// HTTP настройки по умолчанию
		Host:            "0.0.0.0",
		Port:            8080,
		ReadTimeout:     5,
		WriteTimeout:    10,
		IdleTimeout:     120,
		MaxHeaderBytes:  1 << 20, // 1MB
		
		// gRPC настройки по умолчанию
		MaxConcurrentStreams:   1000,
		MaxRecvMsgSize:        4 << 20, // 4MB
		MaxSendMsgSize:        4 << 20, // 4MB
		InitialWindowSize:     65535,
		InitialConnWindowSize: 65535,
		EnableReflection:      true,
		GRPCMiddleware:        make([]grpc.UnaryServerInterceptor, 0),
	}
} 