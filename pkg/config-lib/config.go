package config_lib

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Уровни логирования для внутреннего использования - оставляем для совместимости
const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
)

// Config представляет основную структуру конфигурации
type Config struct {
	// Viper instance
	viper *viper.Viper
	
	// Опции конфигурации
	options Options
	
	// Mutex для безопасности конкурентного доступа
	mu sync.RWMutex
	
	// Флаг, указывающий загружена ли конфигурация
	loaded bool
	
	// Канал для оповещения об изменениях конфигурации
	changes chan struct{}
	
	// Канал для оповещения о закрытии конфигурации
	done chan struct{}
	
	// Колбэк функции для оповещения об изменениях
	onChangeCallbacks []func()
}

// Options опции для создания конфигурации
type Options struct {
	// Имя приложения, используется как префикс для переменных окружения
	AppName string
	
	// Путь к конфигурационному файлу
	ConfigPath string
	
	// Имя конфигурационного файла без расширения
	ConfigName string
	
	// Формат конфигурационного файла (например, "yaml", "json", "toml")
	ConfigType string
	
	// Конфигурация для переопределения через переменные окружения
	EnvPrefix string
	
	// Автоматически заменять точки на подчеркивания в ключах
	EnvKeyReplacer *strings.Replacer
	
	// Автоматический поиск конфигурационного файла
	AutomaticEnv bool
	
	// Наблюдение за изменениями конфигурационного файла
	WatchConfig bool
	
	// Разделитель для вложенных ключей
	KeyDelimiter string
	
	// Значения по умолчанию
	Defaults map[string]interface{}
}

// New создает новый экземпляр Config с указанными опциями
func New(options Options) *Config {
	if options.KeyDelimiter == "" {
		options.KeyDelimiter = "."
	}
	
	// Создаем Viper
	v := viper.New()
	v.SetKeyDelimiter(options.KeyDelimiter)
	
	// Создаем конфигурацию
	config := &Config{
		viper:   v,
		options: options,
		changes: make(chan struct{}, 1),
		done:    make(chan struct{}),
	}
	
	return config
}

// Load загружает конфигурацию
func (c *Config) Load() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	v := c.viper
	opts := c.options
	
	// Устанавливаем имя конфигурационного файла
	if opts.ConfigName != "" {
		v.SetConfigName(opts.ConfigName)
	} else if opts.AppName != "" {
		v.SetConfigName(opts.AppName)
	} else {
		v.SetConfigName("config")
	}
	
	// Устанавливаем путь к конфигурационному файлу
	if opts.ConfigPath != "" {
		v.AddConfigPath(opts.ConfigPath)
	} else {
		// Стандартные пути для поиска конфигурации
		v.AddConfigPath(".")
		v.AddConfigPath("./config")
		v.AddConfigPath("$HOME/.config")
		v.AddConfigPath("/etc")
	}
	
	// Устанавливаем тип конфигурационного файла
	if opts.ConfigType != "" {
		v.SetConfigType(opts.ConfigType)
	}
	
	// Настройка для переменных окружения
	if opts.EnvPrefix != "" {
		v.SetEnvPrefix(opts.EnvPrefix)
	} else if opts.AppName != "" {
		v.SetEnvPrefix(strings.ToUpper(opts.AppName))
	}
	
	// Заменяем точки в ключах на подчеркивания для переменных окружения
	if opts.EnvKeyReplacer != nil {
		v.SetEnvKeyReplacer(opts.EnvKeyReplacer)
	} else {
		v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	}
	
	// Автоматический поиск переменных окружения
	if opts.AutomaticEnv {
		v.AutomaticEnv()
	}
	
	// Устанавливаем значения по умолчанию
	for key, value := range opts.Defaults {
		v.SetDefault(key, value)
	}
	
	// Читаем конфигурационный файл
	err := v.ReadInConfig()
	if err != nil {
		// Если файл не найден, это не всегда ошибка
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("error reading config file: %w", err)
		}
		// Файл не найден, но это нормально - используем значения по умолчанию и переменные окружения
	}
	
	// Если включен режим наблюдения за файлом, настраиваем его
	if opts.WatchConfig {
		v.WatchConfig()
		v.OnConfigChange(func(e fsnotify.Event) {
			select {
			case c.changes <- struct{}{}:
			default:
			}
			
			// Вызываем все колбэки при изменении
			for _, callback := range c.onChangeCallbacks {
				go callback()
			}
		})
	}
	
	c.loaded = true
	return nil
}

// IsLoaded возвращает true, если конфигурация была загружена
func (c *Config) IsLoaded() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.loaded
}

// OnChange регистрирует функцию-колбэк, которая вызывается при изменении конфигурации
func (c *Config) OnChange(callback func()) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.onChangeCallbacks = append(c.onChangeCallbacks, callback)
}

// Close освобождает ресурсы, используемые конфигурацией
func (c *Config) Close() {
	close(c.done)
	close(c.changes)
}

// Changes возвращает канал, сигнализирующий об изменениях конфигурации
func (c *Config) Changes() <-chan struct{} {
	return c.changes
}

// Viper возвращает внутренний экземпляр viper
func (c *Config) Viper() *viper.Viper {
	return c.viper
}

// Get получает значение по ключу
func (c *Config) Get(key string) interface{} {
	return c.viper.Get(key)
}

// GetString получает строковое значение по ключу
func (c *Config) GetString(key string) string {
	return c.viper.GetString(key)
}

// GetBool получает булево значение по ключу
func (c *Config) GetBool(key string) bool {
	return c.viper.GetBool(key)
}

// GetInt получает целочисленное значение по ключу
func (c *Config) GetInt(key string) int {
	return c.viper.GetInt(key)
}

// GetInt64 получает целочисленное значение (int64) по ключу
func (c *Config) GetInt64(key string) int64 {
	return c.viper.GetInt64(key)
}

// GetFloat64 получает значение с плавающей точкой по ключу
func (c *Config) GetFloat64(key string) float64 {
	return c.viper.GetFloat64(key)
}

// GetTime получает временное значение по ключу
func (c *Config) GetTime(key string) time.Time {
	return c.viper.GetTime(key)
}

// GetDuration получает продолжительность по ключу
func (c *Config) GetDuration(key string) time.Duration {
	return c.viper.GetDuration(key)
}

// GetIntSlice получает слайс целых чисел по ключу
func (c *Config) GetIntSlice(key string) []int {
	return c.viper.GetIntSlice(key)
}

// GetStringSlice получает слайс строк по ключу
func (c *Config) GetStringSlice(key string) []string {
	return c.viper.GetStringSlice(key)
}

// GetStringMap получает карту строк по ключу
func (c *Config) GetStringMap(key string) map[string]interface{} {
	return c.viper.GetStringMap(key)
}

// GetStringMapString получает карту строковых значений по ключу
func (c *Config) GetStringMapString(key string) map[string]string {
	return c.viper.GetStringMapString(key)
}

// IsSet проверяет, установлен ли ключ в конфигурации
func (c *Config) IsSet(key string) bool {
	return c.viper.IsSet(key)
}

// Set устанавливает значение в конфигурации
func (c *Config) Set(key string, value interface{}) {
	c.viper.Set(key, value)
}

// SetDefault устанавливает значение по умолчанию
func (c *Config) SetDefault(key string, value interface{}) {
	c.viper.SetDefault(key, value)
}

// GetRequired получает значение по ключу и возвращает ошибку, если оно не установлено
func (c *Config) GetRequired(key string) (interface{}, error) {
	if !c.IsSet(key) {
		return nil, fmt.Errorf("required config key not set: %s", key)
	}
	return c.Get(key), nil
}

// GetRequiredString получает строку по ключу и возвращает ошибку, если она не установлена
func (c *Config) GetRequiredString(key string) (string, error) {
	if !c.IsSet(key) {
		return "", fmt.Errorf("required config key not set: %s", key)
	}
	return c.GetString(key), nil
}

// Unmarshal отображает конфигурацию на структуру
func (c *Config) Unmarshal(rawVal interface{}) error {
	return c.viper.Unmarshal(rawVal)
}

// UnmarshalKey отображает значение ключа на структуру
func (c *Config) UnmarshalKey(key string, rawVal interface{}) error {
	return c.viper.UnmarshalKey(key, rawVal)
}

// ConfigFileUsed возвращает путь к используемому файлу конфигурации
func (c *Config) ConfigFileUsed() string {
	return c.viper.ConfigFileUsed()
}

// LoadAndValidate загружает конфигурацию и проверяет наличие обязательных ключей
func (c *Config) LoadAndValidate(requiredKeys []string) error {
	if err := c.Load(); err != nil {
		return err
	}
	
	if len(requiredKeys) > 0 {
		var missingKeys []string
		for _, key := range requiredKeys {
			if !c.IsSet(key) {
				missingKeys = append(missingKeys, key)
			}
		}
		
		if len(missingKeys) > 0 {
			return fmt.Errorf("missing required config keys: %s", strings.Join(missingKeys, ", "))
		}
	}
	
	return nil
}

// DefaultConfig создает экземпляр конфигурации с настройками по умолчанию
func DefaultConfig(appName string) *Config {
	options := Options{
		AppName:      appName,
		ConfigName:   "config",
		ConfigType:   "yaml",
		AutomaticEnv: true,
		WatchConfig:  true,
		EnvPrefix:    strings.ToUpper(appName),
	}
	
	return New(options)
} 