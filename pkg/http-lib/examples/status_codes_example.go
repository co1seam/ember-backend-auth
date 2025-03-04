package main

import (
	"fmt"
	"log"

	http_lib "github.com/co1seam/tuneflow-backend-auth/pkg/http-lib"
	"github.com/valyala/fasthttp"
)

func main() {
	// Создаем сервер и добавляем middleware
	server := http_lib.New("8080").
		WithLogger().
		WithRecovery().
		WithCORS("*")

	// Маршрут для демонстрации различных HTTP кодов состояния
	server.GET("/status/:code", handleStatusCode)

	// Success статусы
	server.GET("/success/ok", func(ctx *fasthttp.RequestCtx) {
		http_lib.RespondOK(ctx, map[string]interface{}{
			"message": "Success",
			"code":    http_lib.StatusOK,
		})
	})

	server.POST("/success/created", func(ctx *fasthttp.RequestCtx) {
		http_lib.RespondCreated(ctx, map[string]interface{}{
			"message": "Resource created",
			"code":    http_lib.StatusCreated,
		})
	})

	server.POST("/success/accepted", func(ctx *fasthttp.RequestCtx) {
		http_lib.RespondAccepted(ctx, map[string]interface{}{
			"message": "Request accepted for processing",
			"code":    http_lib.StatusAccepted,
		})
	})

	server.DELETE("/success/no-content", func(ctx *fasthttp.RequestCtx) {
		http_lib.RespondNoContent(ctx)
	})

	// Redirection статусы
	server.GET("/redirect/permanent", func(ctx *fasthttp.RequestCtx) {
		http_lib.RespondMovedPermanently(ctx, "/success/ok")
	})

	server.GET("/redirect/temporary", func(ctx *fasthttp.RequestCtx) {
		http_lib.RespondFound(ctx, "/success/ok")
	})

	server.GET("/redirect/see-other", func(ctx *fasthttp.RequestCtx) {
		http_lib.RespondSeeOther(ctx, "/success/ok")
	})

	// Client Error статусы
	server.GET("/error/bad-request", func(ctx *fasthttp.RequestCtx) {
		http_lib.RespondBadRequest(ctx, "Invalid request parameters")
	})

	server.GET("/error/unauthorized", func(ctx *fasthttp.RequestCtx) {
		http_lib.RespondUnauthorized(ctx, "Authentication required")
	})

	server.GET("/error/forbidden", func(ctx *fasthttp.RequestCtx) {
		http_lib.RespondForbidden(ctx, "Insufficient permissions")
	})

	server.GET("/error/not-found", func(ctx *fasthttp.RequestCtx) {
		http_lib.RespondNotFound(ctx, "Resource not found")
	})

	server.GET("/error/method-not-allowed", func(ctx *fasthttp.RequestCtx) {
		http_lib.RespondMethodNotAllowed(ctx, "Method not allowed, use POST")
	})

	server.GET("/error/conflict", func(ctx *fasthttp.RequestCtx) {
		http_lib.RespondConflict(ctx, "Resource already exists")
	})

	server.GET("/error/too-many-requests", func(ctx *fasthttp.RequestCtx) {
		http_lib.RespondTooManyRequests(ctx, "Rate limit exceeded")
	})

	// Server Error статусы
	server.GET("/error/server", func(ctx *fasthttp.RequestCtx) {
		http_lib.RespondInternalServerError(ctx, "Something went wrong")
	})

	server.GET("/error/not-implemented", func(ctx *fasthttp.RequestCtx) {
		http_lib.RespondNotImplemented(ctx, "Feature not implemented yet")
	})

	server.GET("/error/bad-gateway", func(ctx *fasthttp.RequestCtx) {
		http_lib.RespondBadGateway(ctx, "Invalid response from upstream server")
	})

	server.GET("/error/service-unavailable", func(ctx *fasthttp.RequestCtx) {
		http_lib.RespondServiceUnavailable(ctx, "Service is temporarily unavailable")
	})

	// Функциональный стиль с опциями
	server.GET("/advanced/custom-options", func(ctx *fasthttp.RequestCtx) {
		// Добавляем кастомный заголовок и cookie
		cookie := &fasthttp.Cookie{}
		cookie.SetKey("session")
		cookie.SetValue("test-session")
		cookie.SetHTTPOnly(true)
		cookie.SetSecure(true)

		// Используем функциональные опции
		http_lib.RespondOK(ctx, map[string]interface{}{
			"message": "Response with custom options",
		},
			http_lib.WithHeader("X-Custom-Header", "Custom-Value"),
			http_lib.WithCookie(cookie),
		)
	})

	// Маршрут для вызова паники и тестирования RecoveryMiddleware
	server.GET("/panic", func(ctx *fasthttp.RequestCtx) {
		panic("Test panic recovery")
	})

	fmt.Println("Server running on http://localhost:8080")
	fmt.Println("Available routes:")
	fmt.Println("- GET /status/:code - Return response with specified status code")
	fmt.Println("- GET /success/ok - 200 OK response")
	fmt.Println("- POST /success/created - 201 Created response")
	fmt.Println("- POST /success/accepted - 202 Accepted response")
	fmt.Println("- DELETE /success/no-content - 204 No Content response")
	fmt.Println("- GET /redirect/permanent - 301 Moved Permanently response")
	fmt.Println("- GET /redirect/temporary - 302 Found response")
	fmt.Println("- GET /redirect/see-other - 303 See Other response")
	fmt.Println("- GET /error/bad-request - 400 Bad Request response")
	fmt.Println("- GET /error/unauthorized - 401 Unauthorized response")
	fmt.Println("- GET /error/forbidden - 403 Forbidden response")
	fmt.Println("- GET /error/not-found - 404 Not Found response")
	fmt.Println("- GET /error/method-not-allowed - 405 Method Not Allowed response")
	fmt.Println("- GET /error/conflict - 409 Conflict response")
	fmt.Println("- GET /error/too-many-requests - 429 Too Many Requests response")
	fmt.Println("- GET /error/server - 500 Internal Server Error response")
	fmt.Println("- GET /error/not-implemented - 501 Not Implemented response")
	fmt.Println("- GET /error/bad-gateway - 502 Bad Gateway response")
	fmt.Println("- GET /error/service-unavailable - 503 Service Unavailable response")
	fmt.Println("- GET /advanced/custom-options - Response with custom headers and cookies")
	fmt.Println("- GET /panic - Trigger a panic (to test recovery middleware)")

	// Запускаем сервер
	log.Fatal(server.Start())
}

// handleStatusCode обрабатывает запрос с указанным HTTP кодом
func handleStatusCode(ctx *fasthttp.RequestCtx) {
	// Получаем код из параметров
	codeStr := http_lib.GetParam(ctx, "code")
	var code int
	fmt.Sscanf(codeStr, "%d", &code)

	// Проверяем допустимость кода
	if code < 100 || code >= 600 {
		http_lib.RespondBadRequest(ctx, "Invalid status code, must be between 100 and 599")
		return
	}

	// Определяем категорию кода
	category := http_lib.StatusCategory(code)
	categoryText := fmt.Sprintf("This is a %s status code", category)

	// Получаем стандартное описание для кода
	statusText := http_lib.StatusText(code)

	// Формируем ответ в зависимости от категории
	switch {
	case http_lib.IsSuccess(code):
		// 2xx - возвращаем успешный ответ
		http_lib.RespondJSON(ctx, map[string]interface{}{
			"status":   statusText,
			"code":     code,
			"category": categoryText,
		}, http_lib.WithStatusCode(code))

	case http_lib.IsRedirection(code):
		// 3xx - выполняем перенаправление
		ctx.Response.Header.Set("Location", "/success/ok")
		http_lib.RespondStatusCode(ctx, code)

	case http_lib.IsClientError(code), http_lib.IsServerError(code):
		// 4xx, 5xx - возвращаем ошибку
		http_lib.RespondError(ctx, code, statusText)

	default:
		// 1xx и другие - просто устанавливаем код
		http_lib.RespondStatusCode(ctx, code)
	}
} 