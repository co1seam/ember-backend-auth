# Примеры использования HTTP-LIB

В этой директории находятся примеры использования библиотеки HTTP-LIB.

## Запуск примеров

Для запуска примера, используйте следующую команду:

```bash
go run <имя_файла>.go
```

## Примеры

### status_codes_example.go

Демонстрирует использование абстрактной системы HTTP кодов состояния. Запустите пример и перейдите по адресу `http://localhost:8080`.

Доступные маршруты:

- `GET /status/:code` - Возвращает ответ с указанным кодом состояния
- `GET /success/ok` - 200 OK ответ
- `POST /success/created` - 201 Created ответ
- `POST /success/accepted` - 202 Accepted ответ
- `DELETE /success/no-content` - 204 No Content ответ
- `GET /redirect/permanent` - 301 Moved Permanently ответ
- `GET /redirect/temporary` - 302 Found ответ
- `GET /redirect/see-other` - 303 See Other ответ
- `GET /error/bad-request` - 400 Bad Request ответ
- `GET /error/unauthorized` - 401 Unauthorized ответ
- `GET /error/forbidden` - 403 Forbidden ответ
- `GET /error/not-found` - 404 Not Found ответ
- `GET /error/method-not-allowed` - 405 Method Not Allowed ответ
- `GET /error/conflict` - 409 Conflict ответ
- `GET /error/too-many-requests` - 429 Too Many Requests ответ
- `GET /error/server` - 500 Internal Server Error ответ
- `GET /error/not-implemented` - 501 Not Implemented ответ
- `GET /error/bad-gateway` - 502 Bad Gateway ответ
- `GET /error/service-unavailable` - 503 Service Unavailable ответ
- `GET /advanced/custom-options` - Ответ с пользовательскими заголовками и cookies
- `GET /panic` - Вызывает панику (для тестирования middleware восстановления)

## Тестирование с использованием curl

Вы можете использовать следующие команды curl для тестирования примеров:

```bash
# Тестирование различных кодов состояния
curl -i http://localhost:8080/status/200
curl -i http://localhost:8080/status/404
curl -i http://localhost:8080/status/500

# Тестирование успешных ответов
curl -i http://localhost:8080/success/ok
curl -i -X POST http://localhost:8080/success/created
curl -i -X POST http://localhost:8080/success/accepted
curl -i -X DELETE http://localhost:8080/success/no-content

# Тестирование ошибок клиента
curl -i http://localhost:8080/error/bad-request
curl -i http://localhost:8080/error/unauthorized
curl -i http://localhost:8080/error/forbidden
curl -i http://localhost:8080/error/not-found

# Тестирование ошибок сервера
curl -i http://localhost:8080/error/server
curl -i http://localhost:8080/error/not-implemented
curl -i http://localhost:8080/error/service-unavailable

# Тестирование пользовательских опций
curl -i http://localhost:8080/advanced/custom-options

# Тестирование middleware восстановления после паники
curl -i http://localhost:8080/panic
``` 