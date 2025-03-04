package http_lib

import (
	"github.com/valyala/fasthttp"
)

// CreateSimpleServer создает простой сервер с базовыми маршрутами
func CreateSimpleServer(port string) *Server {
	server := New(port)
	
	// Добавляем middleware для логирования
	server.Use(LoggingMiddleware())
	
	// Настраиваем CORS
	server.Use(CORSMiddleware("*"))
	
	// Добавляем обработчик GET запроса
	server.GET("/hello", func(ctx *fasthttp.RequestCtx) {
		SuccessResponse(ctx, map[string]string{
			"message": "Hello, World!",
		})
	})
	
	// Добавляем обработчик POST запроса с JSON
	server.POST("/api/users", func(ctx *fasthttp.RequestCtx) {
		// Структура для разбора JSON
		type UserRequest struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}
		
		var user UserRequest
		if err := ParseJSON(ctx, &user); err != nil {
			BadRequestResponse(ctx, "Invalid JSON")
			return
		}
		
		// Проверяем данные
		if user.Name == "" || user.Email == "" {
			BadRequestResponse(ctx, "Name and email are required")
			return
		}
		
		// Создаем ответ
		CreatedResponse(ctx, map[string]interface{}{
			"id":    123, // В реальном приложении это был бы ID из базы данных
			"name":  user.Name,
			"email": user.Email,
		})
	})
	
	return server
}

// ExampleAPI демонстрирует API для работы с задачами
func ExampleAPI() *Router {
	router := NewRouter()
	
	// Добавляем middleware
	router.Use(LoggingMiddleware())
	router.Use(CORSMiddleware("*"))
	
	// Получение списка задач
	router.GET("/api/tasks", func(ctx *fasthttp.RequestCtx) {
		tasks := []map[string]interface{}{
			{"id": 1, "title": "Task 1", "completed": false},
			{"id": 2, "title": "Task 2", "completed": true},
		}
		
		SuccessResponse(ctx, map[string]interface{}{
			"tasks": tasks,
			"count": len(tasks),
		})
	})
	
	// Получение задачи по ID
	router.GET("/api/tasks/{id}", func(ctx *fasthttp.RequestCtx) {
		id := GetParam(ctx, "id")
		if id == "" {
			BadRequestResponse(ctx, "Task ID is required")
			return
		}
		
		// В реальном приложении здесь был бы поиск в базе данных
		task := map[string]interface{}{
			"id":        id,
			"title":     "Example Task",
			"completed": false,
		}
		
		SuccessResponse(ctx, task)
	})
	
	// Создание новой задачи
	router.POST("/api/tasks", func(ctx *fasthttp.RequestCtx) {
		type TaskRequest struct {
			Title string `json:"title"`
		}
		
		var task TaskRequest
		if err := ParseJSON(ctx, &task); err != nil {
			BadRequestResponse(ctx, "Invalid JSON")
			return
		}
		
		if task.Title == "" {
			BadRequestResponse(ctx, "Title is required")
			return
		}
		
		// В реальном приложении здесь было бы сохранение в базу данных
		CreatedResponse(ctx, map[string]interface{}{
			"id":        123,
			"title":     task.Title,
			"completed": false,
		})
	})
	
	// Обновление задачи
	router.PUT("/api/tasks/{id}", func(ctx *fasthttp.RequestCtx) {
		id := GetParam(ctx, "id")
		if id == "" {
			BadRequestResponse(ctx, "Task ID is required")
			return
		}
		
		type TaskUpdateRequest struct {
			Title     *string `json:"title"`
			Completed *bool   `json:"completed"`
		}
		
		var update TaskUpdateRequest
		if err := ParseJSON(ctx, &update); err != nil {
			BadRequestResponse(ctx, "Invalid JSON")
			return
		}
		
		// В реальном приложении здесь было бы обновление в базе данных
		task := map[string]interface{}{
			"id": id,
		}
		
		if update.Title != nil {
			task["title"] = *update.Title
		} else {
			task["title"] = "Example Task"
		}
		
		if update.Completed != nil {
			task["completed"] = *update.Completed
		} else {
			task["completed"] = false
		}
		
		SuccessResponse(ctx, task)
	})
	
	// Удаление задачи
	router.DELETE("/api/tasks/{id}", func(ctx *fasthttp.RequestCtx) {
		id := GetParam(ctx, "id")
		if id == "" {
			BadRequestResponse(ctx, "Task ID is required")
			return
		}
		
		// В реальном приложении здесь было бы удаление из базы данных
		NoContentResponse(ctx)
	})
	
	return router
} 