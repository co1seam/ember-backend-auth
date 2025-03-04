package config_lib

import (
	"fmt"
	"log"
	"time"
)

// DatabaseConfig пример структуры конфигурации для базы данных
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSL      bool   `mapstructure:"ssl"`
	Timeout  int    `mapstructure:"timeout"`
}

// ServerConfig пример структуры конфигурации для сервера
type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	Debug        bool          `mapstructure:"debug"`
}

// AppConfig пример структуры конфигурации приложения
type AppConfig struct {
	Name        string         `mapstructure:"name"`
	Environment string         `mapstructure:"env"`
	LogLevel    string         `mapstructure:"log_level"`
	Server      ServerConfig   `mapstructure:"server"`
	Database    DatabaseConfig `mapstructure:"database"`
	Features    []string       `mapstructure:"features"`
}

// ExampleBasicUsage показывает базовое использование библиотеки
func ExampleBasicUsage() {
	// Создание конфигурации с настройками по умолчанию
	config := DefaultConfig("myapp")
	
	// Загрузка конфигурации
	if err := config.Load(); err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}
	
	// Получение значений из конфигурации
	host := config.GetString("server.host")
	port := config.GetInt("server.port")
	debug := config.GetBool("server.debug")
	
	fmt.Printf("Server: %s:%d (debug: %v)\n", host, port, debug)
	
	// Получение вложенной структуры
	var dbConfig DatabaseConfig
	if err := config.UnmarshalKey("database", &dbConfig); err != nil {
		fmt.Printf("Error unmarshaling database config: %v\n", err)
		return
	}
	
	fmt.Printf("Database: %s@%s:%d/%s\n", dbConfig.Username, dbConfig.Host, dbConfig.Port, dbConfig.Name)
}

// ExampleTwelveFactorApp показывает использование 12-factor конфигурации
func ExampleTwelveFactorApp() {
	// Создание 12-factor конфигурации
	config := TwelveFactorConfig("myapp")
	
	// Проверка окружения
	if config.IsProduction() {
		fmt.Println("Running in production mode")
	} else if config.IsDevelopment() {
		fmt.Println("Running in development mode")
	}
	
	// Загрузка секрета из переменной окружения или файла
	if err := config.LoadSecret("auth.jwt_secret", "JWT_SECRET", "/run/secrets/jwt_secret"); err != nil {
		log.Printf("Warning: %v", err)
	}
	
	// Получение всей конфигурации как структуры
	var appConfig AppConfig
	if err := config.Unmarshal(&appConfig); err != nil {
		fmt.Printf("Error unmarshaling config: %v\n", err)
		return
	}
	
	fmt.Printf("App: %s (env: %s)\n", appConfig.Name, appConfig.Environment)
}

// ExampleDynamicConfig показывает работу с динамически изменяемой конфигурацией
func ExampleDynamicConfig() {
	// Создание конфигурации с наблюдением за изменениями
	options := Options{
		AppName:     "myapp",
		WatchConfig: true,
	}
	
	config := New(options)
	if err := config.Load(); err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}
	
	// Подписка на изменения конфигурации
	config.OnChange(func() {
		fmt.Println("Config changed!")
		
		// Получаем обновленные значения
		newHost := config.GetString("server.host")
		newPort := config.GetInt("server.port")
		
		fmt.Printf("Updated server: %s:%d\n", newHost, newPort)
	})
	
	// Ожидание изменений (в реальном приложении это было бы в горутине)
	select {
	case <-config.Changes():
		fmt.Println("Received config change notification")
	case <-time.After(1 * time.Minute):
		fmt.Println("No config changes detected within a minute")
	}
}

// ExampleMultipleEnvironments показывает работу с разными окружениями
func ExampleMultipleEnvironments() {
	// Загрузка конфигурации в зависимости от окружения
	config, err := LoadEnvBasedConfig("myapp")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}
	
	// Получение текущего окружения
	env := config.GetEnvironment()
	
	// Разная логика в зависимости от окружения
	switch env {
	case EnvDevelopment:
		fmt.Println("Development mode: enabling detailed logging")
	case EnvTest:
		fmt.Println("Test mode: using mock services")
	case EnvStaging:
		fmt.Println("Staging mode: using staging services")
	case EnvProduction:
		fmt.Println("Production mode: optimizing for performance")
	default:
		fmt.Printf("Unknown environment: %s\n", env)
	}
} 