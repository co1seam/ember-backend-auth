# HTTP-Lib

Гибкая библиотека для работы с HTTP и gRPC в Go, предоставляющая удобный интерфейс для создания серверов.

## Особенности

- **Поддержка HTTP и gRPC** - возможность создавать как HTTP, так и gRPC серверы
- **Гибкая конфигурация** - настраиваемые параметры для обоих типов серверов
- **Middleware поддержка** - для HTTP и gRPC
- **Безопасность** - встроенная поддержка CORS, rate limiting и аутентификации
- **Мониторинг** - встроенные метрики и логирование
- **Удобный API** - простой интерфейс для создания и управления серверами

## Установка

```bash
go get github.com/co1seam/tuneflow-backend-auth/pkg/http-lib
```

## Быстрый старт

### HTTP сервер

```go
package main

import (
	"fmt"
	"log"
	
	http_lib "github.com/co1seam/tuneflow-backend-auth/pkg/http-lib"
)

func main() {
	// Создаем конфигурацию сервера
	config := http_lib.DefaultServerConfig()
	
	// Создаем сервер
	server := http_lib.New(config)
	
	// Добавляем обработчики
	server.GET("/", func(ctx *http_lib.Context) error {
		return ctx.JSON(http_lib.StatusOK, map[string]string{
			"message": "Hello, World!",
		})
	})
	
	// Запускаем сервер
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
```

### gRPC сервер

```go
package main

import (
	"context"
	"log"
	
	http_lib "github.com/co1seam/tuneflow-backend-auth/pkg/http-lib"
	pb "your/proto/package"
)

// server реализует gRPC сервис
type server struct {
	pb.UnimplementedYourServiceServer
}

func (s *server) YourMethod(ctx context.Context, req *pb.YourRequest) (*pb.YourResponse, error) {
	return &pb.YourResponse{
		Message: "Hello from gRPC!",
	}, nil
}

func main() {
	// Создаем конфигурацию сервера
	config := http_lib.DefaultServerConfig()
	config.Port = 50051 // Стандартный порт для gRPC
	
	// Создаем gRPC сервер
	grpcServer := http_lib.NewGRPCServer(config)
	
	// Регистрируем сервис
	pb.RegisterYourServiceServer(grpcServer.GetServer(), &server{})
	
	// Запускаем сервер
	if err := grpcServer.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
```

## Конфигурация

### HTTP сервер

```go
config := http_lib.DefaultServerConfig()
config.Host = "0.0.0.0"
config.Port = 8080
config.ReadTimeout = 5
config.WriteTimeout = 10
config.IdleTimeout = 120
config.MaxHeaderBytes = 1 << 20 // 1MB
```

### gRPC сервер

```go
config := http_lib.DefaultServerConfig()
config.MaxConcurrentStreams = 1000
config.MaxRecvMsgSize = 4 << 20 // 4MB
config.MaxSendMsgSize = 4 << 20 // 4MB
config.InitialWindowSize = 65535
config.InitialConnWindowSize = 65535
config.EnableReflection = true
```

## Middleware

### HTTP Middleware

```go
server := http_lib.New(config)

// Добавляем middleware
server.Use(http_lib.WithCORS())
server.Use(http_lib.WithRateLimit(100, time.Minute))
server.Use(http_lib.WithAuth())
```

### gRPC Middleware

```go
config := http_lib.DefaultServerConfig()
config.GRPCMiddleware = []grpc.UnaryServerInterceptor{
	// Ваш middleware
	func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Логика middleware
		return handler(ctx, req)
	},
}
```

## Примеры

Полные примеры использования можно найти в директории `examples/`:

- `http_example.go` - пример HTTP сервера
- `grpc_example.go` - пример gRPC сервера
- `middleware_example.go` - пример использования middleware
- `status_codes_example.go` - пример использования HTTP статус-кодов

## Лицензия

MIT 