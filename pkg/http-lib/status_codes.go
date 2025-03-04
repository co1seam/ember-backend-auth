package http_lib

// Определение констант для всех стандартных HTTP статус-кодов
const (
	// 2xx - Success
	StatusOK                  = 200
	StatusCreated             = 201
	StatusAccepted            = 202
	StatusNonAuthoritativeInfo = 203
	StatusNoContent           = 204
	StatusResetContent        = 205
	StatusPartialContent      = 206
	StatusMultiStatus         = 207
	StatusAlreadyReported     = 208
	StatusIMUsed              = 226

	// 3xx - Redirection
	StatusMultipleChoices     = 300
	StatusMovedPermanently    = 301
	StatusFound               = 302
	StatusSeeOther            = 303
	StatusNotModified         = 304
	StatusUseProxy            = 305
	StatusTemporaryRedirect   = 307
	StatusPermanentRedirect   = 308

	// 4xx - Client Error
	StatusBadRequest          = 400
	StatusUnauthorized        = 401
	StatusPaymentRequired     = 402
	StatusForbidden           = 403
	StatusNotFound            = 404
	StatusMethodNotAllowed    = 405
	StatusNotAcceptable       = 406
	StatusProxyAuthRequired   = 407
	StatusRequestTimeout      = 408
	StatusConflict            = 409
	StatusGone                = 410
	StatusLengthRequired      = 411
	StatusPreconditionFailed  = 412
	StatusRequestEntityTooLarge = 413
	StatusRequestURITooLong   = 414
	StatusUnsupportedMediaType = 415
	StatusRequestRangeNotSatisfiable = 416
	StatusExpectationFailed   = 417
	StatusTeapot              = 418
	StatusMisdirectedRequest  = 421
	StatusUnprocessableEntity = 422
	StatusLocked              = 423
	StatusFailedDependency    = 424
	StatusTooEarly            = 425
	StatusUpgradeRequired     = 426
	StatusPreconditionRequired = 428
	StatusTooManyRequests     = 429
	StatusRequestHeaderFieldsTooLarge = 431
	StatusUnavailableForLegalReasons = 451

	// 5xx - Server Error
	StatusInternalServerError = 500
	StatusNotImplemented      = 501
	StatusBadGateway          = 502
	StatusServiceUnavailable  = 503
	StatusGatewayTimeout      = 504
	StatusHTTPVersionNotSupported = 505
	StatusVariantAlsoNegotiates = 506
	StatusInsufficientStorage = 507
	StatusLoopDetected        = 508
	StatusNotExtended         = 510
	StatusNetworkAuthenticationRequired = 511
)

// StatusMap содержит соответствие между HTTP кодами и их стандартными описаниями
var StatusMap = map[int]string{
	// 2xx - Success
	StatusOK:                  "OK",
	StatusCreated:             "Created",
	StatusAccepted:            "Accepted",
	StatusNonAuthoritativeInfo: "Non-Authoritative Information",
	StatusNoContent:           "No Content",
	StatusResetContent:        "Reset Content",
	StatusPartialContent:      "Partial Content",
	StatusMultiStatus:         "Multi-Status",
	StatusAlreadyReported:     "Already Reported",
	StatusIMUsed:              "IM Used",

	// 3xx - Redirection
	StatusMultipleChoices:     "Multiple Choices",
	StatusMovedPermanently:    "Moved Permanently",
	StatusFound:               "Found",
	StatusSeeOther:            "See Other",
	StatusNotModified:         "Not Modified",
	StatusUseProxy:            "Use Proxy",
	StatusTemporaryRedirect:   "Temporary Redirect",
	StatusPermanentRedirect:   "Permanent Redirect",

	// 4xx - Client Error
	StatusBadRequest:          "Bad Request",
	StatusUnauthorized:        "Unauthorized",
	StatusPaymentRequired:     "Payment Required",
	StatusForbidden:           "Forbidden",
	StatusNotFound:            "Not Found",
	StatusMethodNotAllowed:    "Method Not Allowed",
	StatusNotAcceptable:       "Not Acceptable",
	StatusProxyAuthRequired:   "Proxy Authentication Required",
	StatusRequestTimeout:      "Request Timeout",
	StatusConflict:            "Conflict",
	StatusGone:                "Gone",
	StatusLengthRequired:      "Length Required",
	StatusPreconditionFailed:  "Precondition Failed",
	StatusRequestEntityTooLarge: "Request Entity Too Large",
	StatusRequestURITooLong:   "Request URI Too Long",
	StatusUnsupportedMediaType: "Unsupported Media Type",
	StatusRequestRangeNotSatisfiable: "Requested Range Not Satisfiable",
	StatusExpectationFailed:   "Expectation Failed",
	StatusTeapot:              "I'm a teapot",
	StatusMisdirectedRequest:  "Misdirected Request",
	StatusUnprocessableEntity: "Unprocessable Entity",
	StatusLocked:              "Locked",
	StatusFailedDependency:    "Failed Dependency",
	StatusTooEarly:            "Too Early",
	StatusUpgradeRequired:     "Upgrade Required",
	StatusPreconditionRequired: "Precondition Required",
	StatusTooManyRequests:     "Too Many Requests",
	StatusRequestHeaderFieldsTooLarge: "Request Header Fields Too Large",
	StatusUnavailableForLegalReasons: "Unavailable For Legal Reasons",

	// 5xx - Server Error
	StatusInternalServerError: "Internal Server Error",
	StatusNotImplemented:      "Not Implemented",
	StatusBadGateway:          "Bad Gateway",
	StatusServiceUnavailable:  "Service Unavailable",
	StatusGatewayTimeout:      "Gateway Timeout",
	StatusHTTPVersionNotSupported: "HTTP Version Not Supported",
	StatusVariantAlsoNegotiates: "Variant Also Negotiates",
	StatusInsufficientStorage: "Insufficient Storage",
	StatusLoopDetected:        "Loop Detected",
	StatusNotExtended:         "Not Extended",
	StatusNetworkAuthenticationRequired: "Network Authentication Required",
}

// StatusCategory возвращает категорию для HTTP кода
func StatusCategory(statusCode int) string {
	switch {
	case statusCode >= 100 && statusCode < 200:
		return "Informational"
	case statusCode >= 200 && statusCode < 300:
		return "Success"
	case statusCode >= 300 && statusCode < 400:
		return "Redirection"
	case statusCode >= 400 && statusCode < 500:
		return "ClientError"
	case statusCode >= 500 && statusCode < 600:
		return "ServerError"
	default:
		return "Unknown"
	}
}

// StatusText возвращает стандартное описание для HTTP кода
func StatusText(statusCode int) string {
	if text, ok := StatusMap[statusCode]; ok {
		return text
	}
	return "Unknown Status Code"
}

// StatusBelongsTo проверяет, принадлежит ли HTTP код определенной категории
func StatusBelongsTo(statusCode int, category string) bool {
	return StatusCategory(statusCode) == category
}

// IsInformational проверяет, является ли HTTP код информационным (1xx)
func IsInformational(statusCode int) bool {
	return statusCode >= 100 && statusCode < 200
}

// IsSuccess проверяет, является ли HTTP код успешным (2xx)
func IsSuccess(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}

// IsRedirection проверяет, является ли HTTP код перенаправлением (3xx)
func IsRedirection(statusCode int) bool {
	return statusCode >= 300 && statusCode < 400
}

// IsClientError проверяет, является ли HTTP код ошибкой клиента (4xx)
func IsClientError(statusCode int) bool {
	return statusCode >= 400 && statusCode < 500
}

// IsServerError проверяет, является ли HTTP код ошибкой сервера (5xx)
func IsServerError(statusCode int) bool {
	return statusCode >= 500 && statusCode < 600
}

// IsError проверяет, является ли HTTP код ошибкой (4xx или 5xx)
func IsError(statusCode int) bool {
	return statusCode >= 400 && statusCode < 600
} 