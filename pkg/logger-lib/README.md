# Logger-Lib

Гибкая библиотека для логирования в Go, основанная на [slog](https://pkg.go.dev/log/slog).

## Особенности

- **Структурированное логирование** - поддержка структурированных данных
- **Множество форматов** - текстовый и JSON формат
- **Ротация логов** - автоматическое управление файлами логов
- **Цветной вывод** - поддержка цветного вывода в консоль
- **Контекст** - поддержка контекста для трейсинга
- **Уровни логирования** - debug, info, warn, error
- **Дополнительная информация** - временные метки, информация о вызове, стек вызовов

## Установка

```bash
go get github.com/co1seam/tuneflow-backend-auth/pkg/logger-lib
```

## Быстрый старт

```go
package main

import (
    logger_lib "github.com/co1seam/tuneflow-backend-auth/pkg/logger-lib"
)

func main() {
    // Создаем настройки логгера
    opts := logger_lib.DefaultOptions()
    opts.LogFile = "logs/app.log"
    opts.Level = logger_lib.LevelDebug
    opts.Format = logger_lib.FormatText
    opts.Colorize = true

    // Создаем логгер
    logger, err := logger_lib.New(opts)
    if err != nil {
        panic(err)
    }

    // Используем логгер
    logger.Info("Application started", "version", "1.0.0")
    logger.Debug("Debug message", "key", "value")
    logger.Warn("Warning message", "error", "connection timeout")
    logger.Error("Error message", "error", "database error", "code", 500)
}
```

## Конфигурация

### Настройки по умолчанию

```go
opts := logger_lib.DefaultOptions()
```

Настройки по умолчанию:
- Level: Info
- Format: Text
- LogFile: "" (stdout)
- MaxSize: 100MB
- MaxBackups: 5
- MaxAge: 30 дней
- AddTimestamp: true
- AddCaller: true
- AddStack: false
- Colorize: true

### Пользовательские настройки

```go
opts := &logger_lib.Options{
    Level:       logger_lib.LevelDebug,
    Format:      logger_lib.FormatJSON,
    LogFile:     "logs/app.log",
    MaxSize:     50 << 20, // 50MB
    MaxBackups:  3,
    MaxAge:      7 * 24 * time.Hour, // 7 дней
    AddTimestamp: true,
    AddCaller:    true,
    AddStack:     true,
    Colorize:     false,
}
```

## Использование

### Базовое логирование

```go
logger.Info("Message", "key", "value")
logger.Debug("Debug message", "key", "value")
logger.Warn("Warning message", "error", "error message")
logger.Error("Error message", "error", "error message", "code", 500)
```

### Логирование с контекстом

```go
ctx := context.Background()
logger.WithContext(ctx).Info("Message with context", "request_id", "abc123")
```

### Логирование с дополнительными полями

```go
logger.With("service", "auth", "version", "1.0.0").Info("Message with fields")
```

### Ротация логов

```go
err := logger_lib.RotateLogFile("logs/app.log", 100<<20, 5, 30*24*time.Hour)
if err != nil {
    logger.Error("Failed to rotate log file", "error", err)
}
```

## Форматы вывода

### Текстовый формат

```
2024-03-03T12:34:56Z [INFO] Message {key=value}
Caller: main.go:123
```

### JSON формат

```json
{
  "time": "2024-03-03T12:34:56Z",
  "level": "INFO",
  "message": "Message",
  "fields": {
    "key": "value"
  },
  "caller": "main.go:123"
}
```

## Примеры

Полные примеры использования можно найти в директории `example/`:

- `main.go` - базовый пример использования
- `context.go` - пример работы с контекстом
- `fields.go` - пример работы с дополнительными полями
- `rotation.go` - пример ротации логов

## Лицензия

MIT 