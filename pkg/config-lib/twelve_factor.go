package config_lib

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Константы для типичных 12-factor настроек
const (
	// Общие настройки приложения
	AppEnv           = "APP_ENV"
	AppPort          = "PORT"
	AppHost          = "HOST"
	AppDebug         = "DEBUG"
	
	// Уровни окружения
	EnvDevelopment   = "development"
	EnvTest          = "test"
	EnvStaging       = "staging"
	EnvProduction    = "production"
)

// TwelveFactorConfig создает конфигурацию, следующую принципам 12-factor app
func TwelveFactorConfig(appName string) *Config {
	options := Options{
		AppName:      appName,
		EnvPrefix:    strings.ToUpper(appName),
		AutomaticEnv: true,
		WatchConfig:  true,
		ConfigName:   "config",
		ConfigType:   "yaml",
	}
	
	// Ищем конфигурацию в стандартных местах
	config := New(options)
	
	// Устанавливаем значения по умолчанию
	setTwelveFactorDefaults(config)
	
	// Загружаем конфигурацию
	config.Load()
	
	return config
}

// setTwelveFactorDefaults устанавливает значения по умолчанию для 12-factor приложения
func setTwelveFactorDefaults(config *Config) {
	config.SetDefault("app.env", EnvDevelopment)
	config.SetDefault("app.port", "8080")
	config.SetDefault("app.host", "0.0.0.0")
	config.SetDefault("app.debug", "false")
}

// LoadEnvBasedConfig загружает конфигурацию в зависимости от окружения (development, test, staging, production)
func LoadEnvBasedConfig(appName string) (*Config, error) {
	// Определяем окружение
	env := os.Getenv(AppEnv)
	if env == "" {
		env = EnvDevelopment
	}
	
	// Имена файлов конфигурации
	configFiles := []string{
		"config",                 // Базовая конфигурация
		fmt.Sprintf("config.%s", env), // Конфигурация для окружения
	}
	
	// Создаем конфигурацию
	options := Options{
		AppName:      appName,
		AutomaticEnv: true,
		EnvPrefix:    strings.ToUpper(appName),
		WatchConfig:  true,
	}
	
	config := New(options)
	
	// Устанавливаем значения по умолчанию
	setTwelveFactorDefaults(config)
	
	// Добавляем стандартные пути поиска
	v := config.viper
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.AddConfigPath("$HOME/.config/" + appName)
	v.AddConfigPath("/etc/" + appName)
	
	// Пытаемся загрузить каждый файл конфигурации
	var configFound bool
	var configErrors []string
	
	for _, name := range configFiles {
		v.SetConfigName(name)
		if err := v.MergeInConfig(); err != nil {
			// Если это не ошибка "файл не найден", добавляем в список ошибок
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				configErrors = append(configErrors, fmt.Sprintf("error loading config '%s': %v", name, err))
			}
		} else {
			configFound = true
		}
	}
	
	// Если есть ошибки, возвращаем их
	if len(configErrors) > 0 {
		return nil, fmt.Errorf("config errors: %s", strings.Join(configErrors, "; "))
	}
	
	// Если ни один конфиг не найден, это не ошибка, но можно сообщить об этом
	// Пользователь сам решит, как обработать эту информацию
	
	// Важная деталь: после загрузки конфигурации устанавливаем флаг loaded
	config.mu.Lock()
	config.loaded = true
	config.mu.Unlock()
	
	return config, nil
}

// GetEnvironment возвращает текущее окружение (development, test, staging, production)
func (c *Config) GetEnvironment() string {
	env := c.GetString("app.env")
	if env == "" {
		env = os.Getenv(AppEnv)
	}
	if env == "" {
		env = EnvDevelopment
	}
	return env
}

// IsDevelopment проверяет, является ли текущее окружение development
func (c *Config) IsDevelopment() bool {
	return c.GetEnvironment() == EnvDevelopment
}

// IsTest проверяет, является ли текущее окружение test
func (c *Config) IsTest() bool {
	return c.GetEnvironment() == EnvTest
}

// IsStaging проверяет, является ли текущее окружение staging
func (c *Config) IsStaging() bool {
	return c.GetEnvironment() == EnvStaging
}

// IsProduction проверяет, является ли текущее окружение production
func (c *Config) IsProduction() bool {
	return c.GetEnvironment() == EnvProduction
}

// LoadSecret загружает секрет из файла или переменной окружения
func (c *Config) LoadSecret(key, envVar, filePath string) error {
	// Сначала проверяем, установлено ли значение
	if c.IsSet(key) {
		return nil
	}
	
	// Проверяем переменную окружения
	if envVar != "" {
		if value := os.Getenv(envVar); value != "" {
			c.Set(key, value)
			return nil
		}
	}
	
	// Проверяем файл
	if filePath != "" {
		data, err := ioutil.ReadFile(filePath)
		if err == nil {
			c.Set(key, strings.TrimSpace(string(data)))
			return nil
		}
		return fmt.Errorf("failed to load secret from file: %w", err)
	}
	
	return fmt.Errorf("secret for key '%s' not found in environment or file", key)
} 