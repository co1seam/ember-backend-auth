package http_lib

import (
	"encoding/json"
	"errors"

	"github.com/valyala/fasthttp"
)

var (
	// ErrInvalidJSON возникает при ошибке разбора JSON
	ErrInvalidJSON = errors.New("invalid JSON")
	
	// ErrMissingParam возникает при отсутствии обязательного параметра
	ErrMissingParam = errors.New("missing required parameter")
)

// ParseJSON разбирает JSON из тела запроса
func ParseJSON(ctx *fasthttp.RequestCtx, dst interface{}) error {
	if !json.Valid(ctx.PostBody()) {
		return ErrInvalidJSON
	}
	
	return json.Unmarshal(ctx.PostBody(), dst)
}

// GetParam получает параметр из URL или query string
func GetParam(ctx *fasthttp.RequestCtx, name string) string {
	// Сначала ищем в URL параметрах
	param := ctx.UserValue(name)
	if param != nil {
		if strParam, ok := param.(string); ok {
			return strParam
		}
	}
	
	// Затем ищем в query string
	return string(ctx.QueryArgs().Peek(name))
}

// GetRequiredParam получает обязательный параметр
func GetRequiredParam(ctx *fasthttp.RequestCtx, name string) (string, error) {
	value := GetParam(ctx, name)
	if value == "" {
		return "", ErrMissingParam
	}
	return value, nil
}

// GetHeader получает значение заголовка
func GetHeader(ctx *fasthttp.RequestCtx, name string) string {
	return string(ctx.Request.Header.Peek(name))
}

// IsJSON проверяет, является ли запрос JSON запросом
func IsJSON(ctx *fasthttp.RequestCtx) bool {
	contentType := string(ctx.Request.Header.Peek("Content-Type"))
	return contentType == "application/json" || contentType == "application/json; charset=utf-8"
}

// GetFileField получает файл из multipart формы
func GetFileField(ctx *fasthttp.RequestCtx, name string) (*fasthttp.FileHeader, error) {
	form, err := ctx.MultipartForm()
	if err != nil {
		return nil, err
	}
	
	fileHeaders := form.File[name]
	if len(fileHeaders) == 0 {
		return nil, ErrMissingParam
	}
	
	return fileHeaders[0], nil
}

// GetFormValue получает значение поля формы
func GetFormValue(ctx *fasthttp.RequestCtx, name string) string {
	return string(ctx.FormValue(name))
} 