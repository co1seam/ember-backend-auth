package http_lib

import (
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
)

// Server представляет HTTP сервер
type Server struct {
	server *fasthttp.Server
	router *Router
	config Config
}

// New создает новый сервер с конфигурацией по умолчанию
func New(port string) *Server {
	config := DefaultConfig()
	config.Port = port
	
	return NewWithConfig(config)
}

// NewWithConfig создает новый сервер с указанной конфигурацией
func NewWithConfig(config Config) *Server {
	router := NewRouter()
	
	server := &Server{
		router: router,
		config: config,
	}
	
	server.server = &fasthttp.Server{
		Addr:           ":" + config.Port,
		ReadTimeout:    config.ReadTimeout,
		WriteTimeout:   config.WriteTimeout,
		MaxHeaderBytes: config.MaxHeaderBytes,
		Handler:        server.handleRequest,
	}
	
	return server
}

// Start запускает сервер
func (s *Server) Start() error {
	fmt.Printf("Server is running at http://localhost:%s\n", s.config.Port)
	return s.server.ListenAndServe()
}

// Stop останавливает сервер
func (s *Server) Stop() error {
	return s.server.Shutdown()
}

// Router возвращает роутер для регистрации маршрутов
func (s *Server) Router() *Router {
	return s.router
}

// HandleFunc добавляет обработчик для указанного метода и пути
func (s *Server) HandleFunc(method, path string, handler fasthttp.RequestHandler) *Server {
	s.router.Handle(method, path, handler)
	return s
}

// GET добавляет GET обработчик
func (s *Server) GET(path string, handler fasthttp.RequestHandler) *Server {
	s.router.GET(path, handler)
	return s
}

// POST добавляет POST обработчик
func (s *Server) POST(path string, handler fasthttp.RequestHandler) *Server {
	s.router.POST(path, handler)
	return s
}

// PUT добавляет PUT обработчик
func (s *Server) PUT(path string, handler fasthttp.RequestHandler) *Server {
	s.router.PUT(path, handler)
	return s
}

// DELETE добавляет DELETE обработчик
func (s *Server) DELETE(path string, handler fasthttp.RequestHandler) *Server {
	s.router.DELETE(path, handler)
	return s
}

// HEAD добавляет HEAD обработчик
func (s *Server) HEAD(path string, handler fasthttp.RequestHandler) *Server {
	s.router.HEAD(path, handler)
	return s
}

// OPTIONS добавляет OPTIONS обработчик
func (s *Server) OPTIONS(path string, handler fasthttp.RequestHandler) *Server {
	s.router.OPTIONS(path, handler)
	return s
}

// PATCH добавляет PATCH обработчик
func (s *Server) PATCH(path string, handler fasthttp.RequestHandler) *Server {
	s.router.PATCH(path, handler)
	return s
}

// Use добавляет middleware
func (s *Server) Use(middleware Middleware) *Server {
	s.router.Use(middleware)
	return s
}

// WithTimeout устанавливает таймауты сервера
func (s *Server) WithTimeout(readTimeout, writeTimeout time.Duration) *Server {
	s.server.ReadTimeout = readTimeout
	s.server.WriteTimeout = writeTimeout
	return s
}

// WithMaxHeaderBytes устанавливает максимальный размер заголовков
func (s *Server) WithMaxHeaderBytes(maxHeaderBytes int) *Server {
	s.server.MaxHeaderBytes = maxHeaderBytes
	return s
}

// WithRecovery добавляет middleware для восстановления после паники
func (s *Server) WithRecovery() *Server {
	s.Use(RecoveryMiddleware())
	return s
}

// WithLogger добавляет middleware для логирования запросов
func (s *Server) WithLogger() *Server {
	s.Use(LoggingMiddleware())
	return s
}

// WithCORS добавляет middleware для настройки CORS
func (s *Server) WithCORS(allowOrigin string) *Server {
	s.Use(CORSMiddleware(allowOrigin))
	return s
}

// WithRateLimit добавляет middleware для ограничения частоты запросов
func (s *Server) WithRateLimit(maxRequestsPerMinute int) *Server {
	s.Use(RateLimitMiddleware(maxRequestsPerMinute))
	return s
}

// WithAuth добавляет middleware для проверки аутентификации
func (s *Server) WithAuth(authFunc func(token string) bool) *Server {
	s.Use(AuthMiddleware(authFunc))
	return s
}

// handleRequest обрабатывает входящие запросы
func (s *Server) handleRequest(ctx *fasthttp.RequestCtx) {
	s.router.Serve(ctx)
}

