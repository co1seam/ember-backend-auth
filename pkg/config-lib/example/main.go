package main

import (
	"fmt"
	"log"
	"time"

	config_lib "github.com/co1seam/tuneflow-backend-auth/pkg/config-lib"
)

// AppConfig представляет конфигурацию приложения
type AppConfig struct {
	Name        string         `mapstructure:"name"`
	Environment string         `mapstructure:"env"`
	LogLevel    string         `mapstructure:"log_level"`
	Version     string         `mapstructure:"version"`
	Server      ServerConfig   `mapstructure:"server"`
	Database    DatabaseConfig `mapstructure:"database"`
	Features    []string       `mapstructure:"features"`
}

// ServerConfig представляет конфигурацию сервера
type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	Debug        bool          `mapstructure:"debug"`
}

// DatabaseConfig представляет конфигурацию базы данных
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSL      bool   `mapstructure:"ssl"`
}

func main() {
	// Создаем конфигурацию с настройками по умолчанию
	config := config_lib.DefaultConfig("example-app")

	// Устанавливаем значения по умолчанию
	config.SetDefault("app.name", "Example App")
	config.SetDefault("app.env", config_lib.EnvDevelopment)
	config.SetDefault("app.log_level", "info")
	config.SetDefault("app.version", "1.0.0")

	config.SetDefault("server.host", "localhost")
	config.SetDefault("server.port", 8080)
	config.SetDefault("server.read_timeout", "5s")
	config.SetDefault("server.write_timeout", "10s")
	config.SetDefault("server.debug", true)

	config.SetDefault("database.host", "localhost")
	config.SetDefault("database.port", 5432)
	config.SetDefault("database.username", "postgres")
	config.SetDefault("database.password", "postgres")
	config.SetDefault("database.name", "example")
	config.SetDefault("database.ssl", false)

	config.SetDefault("features", []string{"auth", "api"})

	// Загружаем конфигурацию
	if err := config.Load(); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Получаем расположение файла конфигурации
	configFile := config.ConfigFileUsed()
	if configFile != "" {
		fmt.Printf("Using config file: %s\n", configFile)
	} else {
		fmt.Println("No config file used, using environment variables and defaults")
	}

	// Проверяем окружение
	env := config.GetEnvironment()
	fmt.Printf("Environment: %s\n", env)

	if config.IsDevelopment() {
		fmt.Println("Running in development mode")
	} else if config.IsProduction() {
		fmt.Println("Running in production mode")
	}

	// Отображаем конфигурацию на структуру
	var appConfig AppConfig
	if err := config.Unmarshal(&appConfig); err != nil {
		log.Fatalf("Error unmarshaling config: %v", err)
	}

	// Выводим конфигурацию
	fmt.Printf("\nApplication Configuration:\n")
	fmt.Printf("  Name: %s\n", appConfig.Name)
	fmt.Printf("  Environment: %s\n", appConfig.Environment)
	fmt.Printf("  Log Level: %s\n", appConfig.LogLevel)
	fmt.Printf("  Version: %s\n", appConfig.Version)

	fmt.Printf("\nServer Configuration:\n")
	fmt.Printf("  Host: %s\n", appConfig.Server.Host)
	fmt.Printf("  Port: %d\n", appConfig.Server.Port)
	fmt.Printf("  Read Timeout: %v\n", appConfig.Server.ReadTimeout)
	fmt.Printf("  Write Timeout: %v\n", appConfig.Server.WriteTimeout)
	fmt.Printf("  Debug Mode: %v\n", appConfig.Server.Debug)

	fmt.Printf("\nDatabase Configuration:\n")
	fmt.Printf("  Host: %s\n", appConfig.Database.Host)
	fmt.Printf("  Port: %d\n", appConfig.Database.Port)
	fmt.Printf("  Username: %s\n", appConfig.Database.Username)
	fmt.Printf("  Password: %s\n", "********") // Не выводим пароль
	fmt.Printf("  Database Name: %s\n", appConfig.Database.Name)
	fmt.Printf("  SSL Enabled: %v\n", appConfig.Database.SSL)

	fmt.Printf("\nEnabled Features: %v\n", appConfig.Features)

	// Загрузка секрета (например, DB_PASSWORD из переменной окружения)
	// С новой версией библиотеки мы можем проверить ошибку
	if err := config.LoadSecret("database.password", "DB_PASSWORD", ""); err != nil {
		fmt.Printf("Note: %v\n", err)
	} else {
		fmt.Println("Successfully loaded database password from environment")
	}

	// Сохраняем конфигурацию в файл (для демонстрации)
	savePath := "config.generated.yaml"
	fmt.Printf("\nSaving current configuration to: %s\n", savePath)
	if err := config.SaveToFile(savePath); err != nil {
		fmt.Printf("Error saving config: %v\n", err)
	} else {
		fmt.Println("Configuration saved successfully!")
	}

	// Демонстрация подписки на изменения (в реальном приложении это было бы в горутине)
	if config.options.WatchConfig {
		fmt.Println("\nWatching for config changes...")
		fmt.Println("Try modifying the config file and see the changes reflected automatically!")

		config.OnChange(func() {
			fmt.Println("\nConfig file changed!")
			config.Unmarshal(&appConfig)
			fmt.Printf("New server port: %d\n", appConfig.Server.Port)
		})

		// Ждем некоторое время для наблюдения за изменениями (в реальном приложении это не требуется)
		fmt.Println("Waiting for 10 seconds to observe config changes...")
		time.Sleep(10 * time.Second)
	}
} 