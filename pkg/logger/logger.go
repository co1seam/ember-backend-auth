package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"time"
)

type Logger struct {
	handler slog.Handler
	ctx     context.Context
}

type Options struct {
	Level     slog.Level
	AddSource bool
	Output    io.Writer
	JSON      bool
}

func New(ctx context.Context, opts Options) *Logger {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	handler := &slog.HandlerOptions{
		Level:     opts.Level,
		AddSource: opts.AddSource,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.String("time", time.Now().UTC().Format(time.RFC3339))
			}
			return a
		},
	}

	logger := slog.Handler(slog.NewJSONHandler(opts.Output, handler))

	return &Logger{
		handler: logger,
		ctx:     ctx,
	}
}

func (l *Logger) With(args ...any) *Logger {
	return &Logger{
		handler: l.handler.WithAttrs(l.toAttrs(args)),
	}
}

func (l *Logger) toAttrs(args []any) []slog.Attr {
	attrs := make([]slog.Attr, 0, len(args)/2)
	for i := 0; i < len(args); i += 2 {
		if i+1 >= len(args) {
			break
		}
		key, val := args[i].(string), args[i+1]
		attrs = append(attrs, slog.Any(key, val))
	}
	return attrs
}

func (l *Logger) log(level slog.Level, msg string, args ...any) error {
	if !l.handler.Enabled(l.ctx, level) {
		return fmt.Errorf("logger isn't enabled")
	}

	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])
	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	r.Add(args...)
	if err := l.handler.Handle(l.ctx, r); err != nil {
		return err
	}

	return nil
}

func (l *Logger) Debug(msg string, args ...any) error {
	return l.log(slog.LevelDebug, msg, args...)
}

func (l *Logger) Info(msg string, args ...any) error {
	return l.log(slog.LevelInfo, msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) error {
	return l.log(slog.LevelWarn, msg, args...)
}

func (l *Logger) Error(msg string, args ...any) error {
	return l.log(slog.LevelError, msg, args...)
}
