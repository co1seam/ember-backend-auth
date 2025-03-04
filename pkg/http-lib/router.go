package http_lib

import (
	"strings"

	"github.com/valyala/fasthttp"
)

// Route представляет маршрут
type Route struct {
	Method  string
	Path    string
	Handler fasthttp.RequestHandler
}

// Router обрабатывает маршруты запросов
type Router struct {
	routes      []Route
	notFound    fasthttp.RequestHandler
	middlewares *MiddlewareChain
}

// NewRouter создает новый роутер
func NewRouter() *Router {
	return &Router{
		routes:      make([]Route, 0),
		middlewares: NewMiddlewareChain(),
		notFound: func(ctx *fasthttp.RequestCtx) {
			RespondNotFound(ctx, "Not Found")
		},
	}
}

// Handle добавляет новый маршрут
func (r *Router) Handle(method, path string, handler fasthttp.RequestHandler) *Router {
	if !IsValidMethod(method) {
		panic("Invalid HTTP method: " + method)
	}
	
	r.routes = append(r.routes, Route{
		Method:  method,
		Path:    path,
		Handler: handler,
	})
	return r
}

// GET добавляет GET маршрут
func (r *Router) GET(path string, handler fasthttp.RequestHandler) *Router {
	return r.Handle(MethodGET, path, handler)
}

// POST добавляет POST маршрут
func (r *Router) POST(path string, handler fasthttp.RequestHandler) *Router {
	return r.Handle(MethodPOST, path, handler)
}

// PUT добавляет PUT маршрут
func (r *Router) PUT(path string, handler fasthttp.RequestHandler) *Router {
	return r.Handle(MethodPUT, path, handler)
}

// DELETE добавляет DELETE маршрут
func (r *Router) DELETE(path string, handler fasthttp.RequestHandler) *Router {
	return r.Handle(MethodDELETE, path, handler)
}

// HEAD добавляет HEAD маршрут
func (r *Router) HEAD(path string, handler fasthttp.RequestHandler) *Router {
	return r.Handle(MethodHEAD, path, handler)
}

// OPTIONS добавляет OPTIONS маршрут
func (r *Router) OPTIONS(path string, handler fasthttp.RequestHandler) *Router {
	return r.Handle(MethodOPTIONS, path, handler)
}

// PATCH добавляет PATCH маршрут
func (r *Router) PATCH(path string, handler fasthttp.RequestHandler) *Router {
	return r.Handle(MethodPATCH, path, handler)
}

// NotFound устанавливает обработчик для 404 ошибки
func (r *Router) NotFound(handler fasthttp.RequestHandler) *Router {
	r.notFound = handler
	return r
}

// Use добавляет middleware в роутер
func (r *Router) Use(middleware Middleware) *Router {
	r.middlewares.Add(middleware)
	return r
}

// Serve обрабатывает HTTP запрос
func (r *Router) Serve(ctx *fasthttp.RequestCtx) {
	// Применяем middleware
	if !r.middlewares.Apply(ctx) {
		return
	}

	// Получаем метод и путь
	method := string(ctx.Method())
	path := string(ctx.Path())

	// Ищем подходящий маршрут
	for _, route := range r.routes {
		if route.Method == method && matchPath(route.Path, path) {
			route.Handler(ctx)
			return
		}
	}

	// Проверяем, является ли запрос OPTIONS и можем ли мы обработать его
	if method == MethodOPTIONS {
		r.handleOptionsRequest(ctx, path)
		return
	}

	// Проверяем, существует ли маршрут с этим путем, но другим методом
	allowedMethods := r.getAllowedMethods(path)
	if len(allowedMethods) > 0 {
		// Возвращаем 405 Method Not Allowed
		ctx.Response.Header.Set("Allow", strings.Join(allowedMethods, ", "))
		RespondMethodNotAllowed(ctx, "Method "+method+" not allowed for this resource")
		return
	}

	// Если маршрут не найден
	r.notFound(ctx)
}

// handleOptionsRequest обрабатывает OPTIONS запрос, возвращая разрешенные методы
func (r *Router) handleOptionsRequest(ctx *fasthttp.RequestCtx, path string) {
	allowedMethods := r.getAllowedMethods(path)
	if len(allowedMethods) > 0 {
		// Возвращаем доступные методы
		ctx.Response.Header.Set("Allow", strings.Join(allowedMethods, ", "))
		RespondStatusCode(ctx, StatusOK)
	} else {
		// Если маршрут не найден, возвращаем 404
		r.notFound(ctx)
	}
}

// getAllowedMethods возвращает список разрешенных методов для указанного пути
func (r *Router) getAllowedMethods(path string) []string {
	var allowedMethods []string
	for _, route := range r.routes {
		if matchPath(route.Path, path) {
			allowedMethods = append(allowedMethods, route.Method)
		}
	}
	return allowedMethods
}

// matchPath проверяет соответствие пути шаблону
// Простая реализация, поддерживает только точное соответствие
func matchPath(pattern, path string) bool {
	// Убираем завершающий слеш, если он есть
	pattern = strings.TrimSuffix(pattern, "/")
	path = strings.TrimSuffix(path, "/")
	
	return pattern == path
} 