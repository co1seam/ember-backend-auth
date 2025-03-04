package main

import (
	"context"
	"time"

	logger_lib "github.com/co1seam/tuneflow-backend-auth/pkg/logger-lib"
)

func main() {
	// Создаем настройки логгера
	opts := logger_lib.DefaultOptions()
	opts.LogFile = "logs/app.log"
	opts.Level = logger_lib.LevelDebug
	opts.Format = logger_lib.FormatText
	opts.Colorize = true
	opts.AddTimestamp = true
	opts.AddCaller = true
	opts.AddStack = true

	// Создаем логгер
	logger, err := logger_lib.New(opts)
	if err != nil {
		panic(err)
	}

	// Примеры использования
	logger.Debug("Debug message", "key", "value")
	logger.Info("Info message", "user_id", 123, "action", "login")
	logger.Warn("Warning message", "error", "connection timeout")
	logger.Error("Error message", "error", "database error", "code", 500)

	// Пример с контекстом
	ctx := context.Background()
	logger.WithContext(ctx).Info("Message with context", "request_id", "abc123")

	// Пример с дополнительными полями
	logger.With("service", "auth", "version", "1.0.0").Info("Message with fields")

	// Пример с ротацией логов
	time.Sleep(time.Second) // Ждем, чтобы увидеть ротацию
	if err := logger_lib.RotateLogFile(opts.LogFile, opts.MaxSize, opts.MaxBackups, opts.MaxAge); err != nil {
		logger.Error("Failed to rotate log file", "error", err)
	}
} 