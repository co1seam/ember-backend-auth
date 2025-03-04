package http_lib

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

// ResponseOption функциональная опция для настройки ответа
type ResponseOption func(*ResponseOptions)

// ResponseOptions содержит опции для настройки ответа
type ResponseOptions struct {
	Headers    map[string]string
	Cookies    []*fasthttp.Cookie
	StatusCode int
}

// DefaultResponseOptions возвращает опции по умолчанию
func DefaultResponseOptions() ResponseOptions {
	return ResponseOptions{
		Headers:    make(map[string]string),
		Cookies:    make([]*fasthttp.Cookie, 0),
		StatusCode: StatusOK,
	}
}

// WithHeader добавляет заголовок в ответ
func WithHeader(name, value string) ResponseOption {
	return func(options *ResponseOptions) {
		options.Headers[name] = value
	}
}

// WithCookie добавляет cookie в ответ
func WithCookie(cookie *fasthttp.Cookie) ResponseOption {
	return func(options *ResponseOptions) {
		options.Cookies = append(options.Cookies, cookie)
	}
}

// WithStatusCode устанавливает HTTP код ответа
func WithStatusCode(statusCode int) ResponseOption {
	return func(options *ResponseOptions) {
		options.StatusCode = statusCode
	}
}

// WithContentType устанавливает Content-Type заголовок
func WithContentType(contentType string) ResponseOption {
	return func(options *ResponseOptions) {
		options.Headers["Content-Type"] = contentType
	}
}

// applyResponseOptions применяет опции к ответу
func applyResponseOptions(ctx *fasthttp.RequestCtx, options ResponseOptions) {
	// Устанавливаем статус код
	ctx.Response.SetStatusCode(options.StatusCode)
	
	// Устанавливаем заголовки
	for name, value := range options.Headers {
		ctx.Response.Header.Set(name, value)
	}
	
	// Устанавливаем cookies
	for _, cookie := range options.Cookies {
		ctx.Response.Header.SetCookie(cookie)
	}
}

// Respond отправляет HTTP ответ с заданными опциями
func Respond(ctx *fasthttp.RequestCtx, body []byte, opts ...ResponseOption) {
	options := DefaultResponseOptions()
	
	// Применяем все опции
	for _, opt := range opts {
		opt(&options)
	}
	
	// Применяем опции к ответу
	applyResponseOptions(ctx, options)
	
	// Устанавливаем тело ответа
	ctx.Response.SetBody(body)
}

// RespondString отправляет строковый HTTP ответ
func RespondString(ctx *fasthttp.RequestCtx, body string, opts ...ResponseOption) {
	// По умолчанию текстовый Content-Type, если не указан иной
	hasContentType := false
	options := DefaultResponseOptions()
	
	for _, opt := range opts {
		opt(&options)
		if _, ok := options.Headers["Content-Type"]; ok {
			hasContentType = true
		}
	}
	
	if !hasContentType {
		opts = append(opts, WithContentType("text/plain; charset=utf-8"))
	}
	
	Respond(ctx, []byte(body), opts...)
}

// RespondJSON отправляет JSON HTTP ответ
func RespondJSON(ctx *fasthttp.RequestCtx, data interface{}, opts ...ResponseOption) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	
	// Добавляем JSON Content-Type, если не указан иной
	hasContentType := false
	options := DefaultResponseOptions()
	
	for _, opt := range opts {
		opt(&options)
		if _, ok := options.Headers["Content-Type"]; ok {
			hasContentType = true
		}
	}
	
	if !hasContentType {
		opts = append(opts, WithContentType("application/json; charset=utf-8"))
	}
	
	Respond(ctx, jsonData, opts...)
	return nil
}

// RespondError отправляет JSON ответ с ошибкой
func RespondError(ctx *fasthttp.RequestCtx, statusCode int, message string, opts ...ResponseOption) error {
	// Устанавливаем статус код и добавляем его в опции
	opts = append(opts, WithStatusCode(statusCode))
	
	// Формируем JSON с ошибкой
	errorData := map[string]interface{}{
		"error": map[string]interface{}{
			"code":    statusCode,
			"message": message,
			"status":  StatusText(statusCode),
		},
	}
	
	return RespondJSON(ctx, errorData, opts...)
}

// RespondStatusCode отправляет ответ только с HTTP кодом без тела
func RespondStatusCode(ctx *fasthttp.RequestCtx, statusCode int, opts ...ResponseOption) {
	opts = append(opts, WithStatusCode(statusCode))
	
	options := DefaultResponseOptions()
	for _, opt := range opts {
		opt(&options)
	}
	
	applyResponseOptions(ctx, options)
}

// ========================= ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ ДЛЯ 2XX =========================

// RespondOK отправляет ответ с кодом 200 OK
func RespondOK(ctx *fasthttp.RequestCtx, data interface{}, opts ...ResponseOption) error {
	opts = append(opts, WithStatusCode(StatusOK))
	return RespondJSON(ctx, data, opts...)
}

// RespondCreated отправляет ответ с кодом 201 Created
func RespondCreated(ctx *fasthttp.RequestCtx, data interface{}, opts ...ResponseOption) error {
	opts = append(opts, WithStatusCode(StatusCreated))
	return RespondJSON(ctx, data, opts...)
}

// RespondAccepted отправляет ответ с кодом 202 Accepted
func RespondAccepted(ctx *fasthttp.RequestCtx, data interface{}, opts ...ResponseOption) error {
	opts = append(opts, WithStatusCode(StatusAccepted))
	return RespondJSON(ctx, data, opts...)
}

// RespondNoContent отправляет ответ с кодом 204 No Content
func RespondNoContent(ctx *fasthttp.RequestCtx, opts ...ResponseOption) {
	opts = append(opts, WithStatusCode(StatusNoContent))
	RespondStatusCode(ctx, StatusNoContent, opts...)
}

// ========================= ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ ДЛЯ 3XX =========================

// RespondMovedPermanently отправляет ответ с кодом 301 Moved Permanently
func RespondMovedPermanently(ctx *fasthttp.RequestCtx, location string, opts ...ResponseOption) {
	opts = append(opts, WithStatusCode(StatusMovedPermanently), WithHeader("Location", location))
	RespondStatusCode(ctx, StatusMovedPermanently, opts...)
}

// RespondFound отправляет ответ с кодом 302 Found
func RespondFound(ctx *fasthttp.RequestCtx, location string, opts ...ResponseOption) {
	opts = append(opts, WithStatusCode(StatusFound), WithHeader("Location", location))
	RespondStatusCode(ctx, StatusFound, opts...)
}

// RespondSeeOther отправляет ответ с кодом 303 See Other
func RespondSeeOther(ctx *fasthttp.RequestCtx, location string, opts ...ResponseOption) {
	opts = append(opts, WithStatusCode(StatusSeeOther), WithHeader("Location", location))
	RespondStatusCode(ctx, StatusSeeOther, opts...)
}

// RespondNotModified отправляет ответ с кодом 304 Not Modified
func RespondNotModified(ctx *fasthttp.RequestCtx, opts ...ResponseOption) {
	opts = append(opts, WithStatusCode(StatusNotModified))
	RespondStatusCode(ctx, StatusNotModified, opts...)
}

// RespondTemporaryRedirect отправляет ответ с кодом 307 Temporary Redirect
func RespondTemporaryRedirect(ctx *fasthttp.RequestCtx, location string, opts ...ResponseOption) {
	opts = append(opts, WithStatusCode(StatusTemporaryRedirect), WithHeader("Location", location))
	RespondStatusCode(ctx, StatusTemporaryRedirect, opts...)
}

// RespondPermanentRedirect отправляет ответ с кодом 308 Permanent Redirect
func RespondPermanentRedirect(ctx *fasthttp.RequestCtx, location string, opts ...ResponseOption) {
	opts = append(opts, WithStatusCode(StatusPermanentRedirect), WithHeader("Location", location))
	RespondStatusCode(ctx, StatusPermanentRedirect, opts...)
}

// ========================= ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ ДЛЯ 4XX =========================

// RespondBadRequest отправляет ответ с кодом 400 Bad Request
func RespondBadRequest(ctx *fasthttp.RequestCtx, message string, opts ...ResponseOption) error {
	if message == "" {
		message = "Bad Request"
	}
	return RespondError(ctx, StatusBadRequest, message, opts...)
}

// RespondUnauthorized отправляет ответ с кодом 401 Unauthorized
func RespondUnauthorized(ctx *fasthttp.RequestCtx, message string, opts ...ResponseOption) error {
	if message == "" {
		message = "Unauthorized"
	}
	return RespondError(ctx, StatusUnauthorized, message, opts...)
}

// RespondForbidden отправляет ответ с кодом 403 Forbidden
func RespondForbidden(ctx *fasthttp.RequestCtx, message string, opts ...ResponseOption) error {
	if message == "" {
		message = "Forbidden"
	}
	return RespondError(ctx, StatusForbidden, message, opts...)
}

// RespondNotFound отправляет ответ с кодом 404 Not Found
func RespondNotFound(ctx *fasthttp.RequestCtx, message string, opts ...ResponseOption) error {
	if message == "" {
		message = "Not Found"
	}
	return RespondError(ctx, StatusNotFound, message, opts...)
}

// RespondMethodNotAllowed отправляет ответ с кодом 405 Method Not Allowed
func RespondMethodNotAllowed(ctx *fasthttp.RequestCtx, message string, opts ...ResponseOption) error {
	if message == "" {
		message = "Method Not Allowed"
	}
	return RespondError(ctx, StatusMethodNotAllowed, message, opts...)
}

// RespondConflict отправляет ответ с кодом 409 Conflict
func RespondConflict(ctx *fasthttp.RequestCtx, message string, opts ...ResponseOption) error {
	if message == "" {
		message = "Conflict"
	}
	return RespondError(ctx, StatusConflict, message, opts...)
}

// RespondUnprocessableEntity отправляет ответ с кодом 422 Unprocessable Entity
func RespondUnprocessableEntity(ctx *fasthttp.RequestCtx, message string, opts ...ResponseOption) error {
	if message == "" {
		message = "Unprocessable Entity"
	}
	return RespondError(ctx, StatusUnprocessableEntity, message, opts...)
}

// RespondTooManyRequests отправляет ответ с кодом 429 Too Many Requests
func RespondTooManyRequests(ctx *fasthttp.RequestCtx, message string, opts ...ResponseOption) error {
	if message == "" {
		message = "Too Many Requests"
	}
	return RespondError(ctx, StatusTooManyRequests, message, opts...)
}

// ========================= ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ ДЛЯ 5XX =========================

// RespondInternalServerError отправляет ответ с кодом 500 Internal Server Error
func RespondInternalServerError(ctx *fasthttp.RequestCtx, message string, opts ...ResponseOption) error {
	if message == "" {
		message = "Internal Server Error"
	}
	return RespondError(ctx, StatusInternalServerError, message, opts...)
}

// RespondNotImplemented отправляет ответ с кодом 501 Not Implemented
func RespondNotImplemented(ctx *fasthttp.RequestCtx, message string, opts ...ResponseOption) error {
	if message == "" {
		message = "Not Implemented"
	}
	return RespondError(ctx, StatusNotImplemented, message, opts...)
}

// RespondBadGateway отправляет ответ с кодом 502 Bad Gateway
func RespondBadGateway(ctx *fasthttp.RequestCtx, message string, opts ...ResponseOption) error {
	if message == "" {
		message = "Bad Gateway"
	}
	return RespondError(ctx, StatusBadGateway, message, opts...)
}

// RespondServiceUnavailable отправляет ответ с кодом 503 Service Unavailable
func RespondServiceUnavailable(ctx *fasthttp.RequestCtx, message string, opts ...ResponseOption) error {
	if message == "" {
		message = "Service Unavailable"
	}
	return RespondError(ctx, StatusServiceUnavailable, message, opts...)
}

// ========================= ПОДДЕРЖКА ПРЕДЫДУЩЕГО API =========================

// Эти функции оставлены для обратной совместимости

// JSONResponse отправляет JSON ответ
func JSONResponse(ctx *fasthttp.RequestCtx, statusCode int, data interface{}) error {
	return RespondJSON(ctx, data, WithStatusCode(statusCode))
}

// ErrorResponse отправляет JSON ответ с ошибкой
func ErrorResponse(ctx *fasthttp.RequestCtx, statusCode int, message string) error {
	return RespondError(ctx, statusCode, message)
}

// SuccessResponse отправляет успешный JSON ответ
func SuccessResponse(ctx *fasthttp.RequestCtx, data interface{}) error {
	return RespondOK(ctx, data)
}

// CreatedResponse отправляет JSON ответ с кодом 201 (Created)
func CreatedResponse(ctx *fasthttp.RequestCtx, data interface{}) error {
	return RespondCreated(ctx, data)
}

// NoContentResponse отправляет ответ с кодом 204 (No Content)
func NoContentResponse(ctx *fasthttp.RequestCtx) error {
	RespondNoContent(ctx)
	return nil
}

// BadRequestResponse отправляет ответ с кодом 400 (Bad Request)
func BadRequestResponse(ctx *fasthttp.RequestCtx, message string) error {
	return RespondBadRequest(ctx, message)
}

// UnauthorizedResponse отправляет ответ с кодом 401 (Unauthorized)
func UnauthorizedResponse(ctx *fasthttp.RequestCtx, message string) error {
	return RespondUnauthorized(ctx, message)
}

// ForbiddenResponse отправляет ответ с кодом 403 (Forbidden)
func ForbiddenResponse(ctx *fasthttp.RequestCtx, message string) error {
	return RespondForbidden(ctx, message)
}

// NotFoundResponse отправляет ответ с кодом 404 (Not Found)
func NotFoundResponse(ctx *fasthttp.RequestCtx, message string) error {
	return RespondNotFound(ctx, message)
}

// InternalServerErrorResponse отправляет ответ с кодом 500 (Internal Server Error)
func InternalServerErrorResponse(ctx *fasthttp.RequestCtx, message string) error {
	return RespondInternalServerError(ctx, message)
} 