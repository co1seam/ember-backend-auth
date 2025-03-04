package http_lib

// Определение констант для всех стандартных HTTP методов
const (
	MethodGET     = "GET"
	MethodPOST    = "POST"
	MethodPUT     = "PUT"
	MethodDELETE  = "DELETE"
	MethodHEAD    = "HEAD"
	MethodOPTIONS = "OPTIONS"
	MethodPATCH   = "PATCH"
	MethodTRACE   = "TRACE"
	MethodCONNECT = "CONNECT"
)

// AllHTTPMethods содержит список всех поддерживаемых HTTP методов
var AllHTTPMethods = []string{
	MethodGET,
	MethodPOST,
	MethodPUT,
	MethodDELETE,
	MethodHEAD,
	MethodOPTIONS,
	MethodPATCH,
	MethodTRACE,
	MethodCONNECT,
}

// SafeHTTPMethods содержит список "безопасных" HTTP методов (не изменяют состояние)
var SafeHTTPMethods = []string{
	MethodGET,
	MethodHEAD,
	MethodOPTIONS,
	MethodTRACE,
}

// IdempotentHTTPMethods содержит список идемпотентных HTTP методов
var IdempotentHTTPMethods = []string{
	MethodGET,
	MethodHEAD,
	MethodPUT,
	MethodDELETE,
	MethodOPTIONS,
	MethodTRACE,
}

// MethodWithRequestBody проверяет, имеет ли HTTP метод тело запроса
func MethodWithRequestBody(method string) bool {
	switch method {
	case MethodPOST, MethodPUT, MethodPATCH:
		return true
	default:
		return false
	}
}

// IsValidMethod проверяет, является ли строка допустимым HTTP методом
func IsValidMethod(method string) bool {
	for _, m := range AllHTTPMethods {
		if m == method {
			return true
		}
	}
	return false
}

// IsSafeMethod проверяет, является ли HTTP метод "безопасным"
func IsSafeMethod(method string) bool {
	for _, m := range SafeHTTPMethods {
		if m == method {
			return true
		}
	}
	return false
}

// IsIdempotentMethod проверяет, является ли HTTP метод идемпотентным
func IsIdempotentMethod(method string) bool {
	for _, m := range IdempotentHTTPMethods {
		if m == method {
			return true
		}
	}
	return false
} 