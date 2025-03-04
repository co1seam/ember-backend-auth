package http_lib

import (
	"github.com/valyala/fasthttp"
)

// Middleware функция для обработки HTTP запросов
type Middleware func(ctx *fasthttp.RequestCtx) error

// MiddlewareChain представляет цепочку middleware
type MiddlewareChain struct {
	middlewares []Middleware
}

// NewMiddlewareChain создает новую цепочку middleware
func NewMiddlewareChain() *MiddlewareChain {
	return &MiddlewareChain{
		middlewares: make([]Middleware, 0),
	}
}

// Add добавляет middleware в цепочку
func (mc *MiddlewareChain) Add(middleware Middleware) *MiddlewareChain {
	mc.middlewares = append(mc.middlewares, middleware)
	return mc
}

// Apply применяет все middleware к запросу
func (mc *MiddlewareChain) Apply(ctx *fasthttp.RequestCtx) bool {
	for _, middleware := range mc.middlewares {
		if err := middleware(ctx); err != nil {
			return false
		}
		// Если статус код уже установлен, прерываем цепочку
		if ctx.Response.StatusCode() != 0 {
			return false
		}
	}
	return true
}

// LoggingMiddleware создает middleware для логирования запросов
func LoggingMiddleware() Middleware {
	return func(ctx *fasthttp.RequestCtx) error {
		method := ctx.Method()
		path := ctx.Path()
		
		// Логируем информацию о запросе (в данной реализации просто выводим в fasthttp Logger)
		ctx.Logger().Printf("Request: %s %s", method, path)
		return nil
	}
}

// CORSMiddleware создает middleware для настройки CORS
func CORSMiddleware(allowOrigin string) Middleware {
	return func(ctx *fasthttp.RequestCtx) error {
		ctx.Response.Header.Set("Access-Control-Allow-Origin", allowOrigin)
		ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		// Если это preflight запрос, завершаем обработку
		if string(ctx.Method()) == MethodOPTIONS {
			RespondStatusCode(ctx, StatusOK)
		}
		
		return nil
	}
}

// RecoveryMiddleware создает middleware для восстановления после паники
func RecoveryMiddleware() Middleware {
	return func(ctx *fasthttp.RequestCtx) error {
		defer func() {
			if err := recover(); err != nil {
				// Логируем ошибку
				ctx.Logger().Printf("Panic recovered: %v", err)
				
				// Отправляем ответ с ошибкой
				RespondInternalServerError(ctx, "Internal Server Error")
			}
		}()
		
		return nil
	}
}

// TimeoutMiddleware создает middleware для проверки таймаута
func TimeoutMiddleware() Middleware {
	return func(ctx *fasthttp.RequestCtx) error {
		if ctx.ConnRequestNum() > 100 {
			RespondServiceUnavailable(ctx, "Service is under high load, please try again later")
			return nil
		}
		
		return nil
	}
}

// RateLimitMiddleware создает middleware для ограничения частоты запросов
// Простая реализация, в реальном приложении нужно использовать счетчики или Redis
func RateLimitMiddleware(maxRequestsPerMinute int) Middleware {
	// Здесь должна быть логика для хранения счетчиков по IP
	return func(ctx *fasthttp.RequestCtx) error {
		// Заглушка для примера
		if maxRequestsPerMinute == 0 {
			RespondTooManyRequests(ctx, "Rate limit exceeded")
			return nil
		}
		
		return nil
	}
}

// AuthMiddleware создает middleware для проверки аутентификации
func AuthMiddleware(authFunc func(token string) bool) Middleware {
	return func(ctx *fasthttp.RequestCtx) error {
		token := string(ctx.Request.Header.Peek("Authorization"))
		
		if token == "" {
			RespondUnauthorized(ctx, "Authorization header is required")
			return nil
		}
		
		// Убираем префикс "Bearer " если он есть
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}
		
		if !authFunc(token) {
			RespondUnauthorized(ctx, "Invalid token")
			return nil
		}
		
		return nil
	}
} 