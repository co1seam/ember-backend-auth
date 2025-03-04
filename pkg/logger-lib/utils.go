package logger_lib

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// LogEntry представляет запись лога
type LogEntry struct {
	Time    time.Time              `json:"time"`
	Level   string                 `json:"level"`
	Message string                 `json:"message"`
	Fields  map[string]interface{} `json:"fields,omitempty"`
	Caller  string                 `json:"caller,omitempty"`
	Stack   string                 `json:"stack,omitempty"`
}

// RotateLogFile ротирует файл логов
func RotateLogFile(filePath string, maxSize int64, maxBackups int, maxAge time.Duration) error {
	// Проверяем размер файла
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if info.Size() < maxSize {
		return nil
	}

	// Создаем имя для нового файла
	dir := filepath.Dir(filePath)
	ext := filepath.Ext(filePath)
	name := filepath.Base(filePath[:len(filePath)-len(ext)])
	timestamp := time.Now().Format("20060102-150405")
	newPath := filepath.Join(dir, fmt.Sprintf("%s-%s%s", name, timestamp, ext))

	// Переименовываем текущий файл
	if err := os.Rename(filePath, newPath); err != nil {
		return err
	}

	// Удаляем старые файлы
	if err := cleanOldLogs(dir, name, maxBackups, maxAge); err != nil {
		return err
	}

	return nil
}

// cleanOldLogs удаляет старые файлы логов
func cleanOldLogs(dir, name string, maxBackups int, maxAge time.Duration) error {
	// Получаем список файлов
	files, err := filepath.Glob(filepath.Join(dir, name+"-*"))
	if err != nil {
		return err
	}

	// Сортируем файлы по времени создания
	var logFiles []struct {
		path    string
		modTime time.Time
	}

	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}

		logFiles = append(logFiles, struct {
			path    string
			modTime time.Time
		}{
			path:    file,
			modTime: info.ModTime(),
		})
	}

	// Сортируем по времени модификации (новые первыми)
	sort.Slice(logFiles, func(i, j int) bool {
		return logFiles[i].modTime.After(logFiles[j].modTime)
	})

	// Удаляем лишние файлы
	for i, file := range logFiles {
		if i >= maxBackups {
			os.Remove(file.path)
			continue
		}

		if time.Since(file.modTime) > maxAge {
			os.Remove(file.path)
		}
	}

	return nil
}

// FormatLogEntry форматирует запись лога
func FormatLogEntry(entry LogEntry, format Format) (string, error) {
	switch format {
	case FormatJSON:
		data, err := json.Marshal(entry)
		if err != nil {
			return "", err
		}
		return string(data), nil
	default:
		return fmt.Sprintf("%s [%s] %s %s %s",
			entry.Time.Format(time.RFC3339),
			entry.Level,
			entry.Message,
			formatFields(entry.Fields),
			formatCaller(entry.Caller, entry.Stack),
		), nil
	}
}

// formatFields форматирует поля лога
func formatFields(fields map[string]interface{}) string {
	if len(fields) == 0 {
		return ""
	}

	var parts []string
	for k, v := range fields {
		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
	}
	return "{" + strings.Join(parts, " ") + "}"
}

// formatCaller форматирует информацию о вызове
func formatCaller(caller, stack string) string {
	if caller == "" {
		return ""
	}

	if stack != "" {
		return fmt.Sprintf("\nCaller: %s\nStack:\n%s", caller, stack)
	}

	return fmt.Sprintf("\nCaller: %s", caller)
} 