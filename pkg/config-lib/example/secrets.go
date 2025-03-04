package main

import (
	"fmt"
	"log"

	config_lib "github.com/co1seam/tuneflow-backend-auth/pkg/config-lib"
)

func main() {
	// Создаем конфигурацию с чувствительными данными
	config := config_lib.DefaultConfig("secrets-example")
	
	// Устанавливаем секретные и несекретные данные
	config.Set("app.name", "Secret App")
	config.Set("app.env", "development")
	
	// Чувствительные данные
	config.Set("database.username", "admin")
	config.Set("database.password", "super-secret-password")
	config.Set("auth.api_key", "12345-very-secret-api-key")
	config.Set("auth.jwt_secret", "jwt-secret-for-signing-tokens")
	
	// Обычные данные
	config.Set("server.host", "localhost")
	config.Set("server.port", 8080)
	
	// Создаем менеджер секретов
	secretsManager := config_lib.NewSecretsManager(config)
	
	// Пароль для шифрования
	encryptionPassphrase := "my-secure-passphrase"
	
	// Сохраняем зашифрованную конфигурацию
	encryptedConfigPath := "config.encrypted.yaml"
	fmt.Printf("Сохраняем зашифрованную конфигурацию в %s\n", encryptedConfigPath)
	if err := secretsManager.CreateEncryptedConfig(encryptedConfigPath, encryptionPassphrase); err != nil {
		log.Fatalf("Ошибка создания зашифрованной конфигурации: %v", err)
	}
	
	// Создаем новую пустую конфигурацию для демонстрации загрузки
	newConfig := config_lib.DefaultConfig("secrets-example")
	newSecretsManager := config_lib.NewSecretsManager(newConfig)
	
	// Загружаем зашифрованную конфигурацию
	fmt.Printf("\nЗагружаем зашифрованную конфигурацию из %s\n", encryptedConfigPath)
	if err := newSecretsManager.LoadEncryptedConfig(encryptedConfigPath, encryptionPassphrase); err != nil {
		log.Fatalf("Ошибка загрузки зашифрованной конфигурации: %v", err)
	}
	
	// Проверяем, что секретные данные были успешно загружены и расшифрованы
	fmt.Println("\nЗагруженная и расшифрованная конфигурация:")
	fmt.Printf("app.name = %s\n", newConfig.GetString("app.name"))
	fmt.Printf("app.env = %s\n", newConfig.GetString("app.env"))
	fmt.Printf("server.host = %s\n", newConfig.GetString("server.host"))
	fmt.Printf("server.port = %d\n", newConfig.GetInt("server.port"))
	
	// Секретные данные
	fmt.Printf("database.username = %s\n", newConfig.GetString("database.username"))
	fmt.Printf("database.password = %s\n", newConfig.GetString("database.password"))
	fmt.Printf("auth.api_key = %s\n", newConfig.GetString("auth.api_key"))
	fmt.Printf("auth.jwt_secret = %s\n", newConfig.GetString("auth.jwt_secret"))
	
	// Демонстрация загрузки секретов из переменных окружения
	fmt.Println("\nДемонстрация загрузки секретов из переменных окружения:")
	
	// Эти команды нужно выполнить в оболочке перед запуском программы:
	// export DB_PASSWORD=env-database-password
	// export API_KEY=env-api-key
	
	if loaded := newSecretsManager.LoadSecretFromEnv("database.password", "DB_PASSWORD"); loaded {
		fmt.Println("Успешно загружен секрет database.password из переменной окружения")
	} else {
		fmt.Println("Секрет database.password не найден в переменной окружения")
	}
	
	if loaded := newSecretsManager.LoadSecretFromEnv("auth.api_key", "API_KEY"); loaded {
		fmt.Println("Успешно загружен секрет auth.api_key из переменной окружения")
	} else {
		fmt.Println("Секрет auth.api_key не найден в переменной окружения")
	}
	
	fmt.Printf("database.password после загрузки из переменной окружения = %s\n", newConfig.GetString("database.password"))
	fmt.Printf("auth.api_key после загрузки из переменной окружения = %s\n", newConfig.GetString("auth.api_key"))
	
	// Демонстрация загрузки секретов из файлов
	fmt.Println("\nПример загрузки секретов из файлов:")
	
	// Файл не существует - продемонстрируем обработку ошибки
	if err := newSecretsManager.LoadSecretFromFile("database.password", "/path/to/nonexistent/db_password.txt"); err != nil {
		fmt.Printf("Ошибка загрузки из файла: %v\n", err)
	}
	
	// Docker и Kubernetes секреты (демонстрация API)
	fmt.Println("\nПримеры других методов (не выполняются, только демонстрация API):")
	fmt.Println("newSecretsManager.LoadDockerSecret(\"auth.jwt_secret\", \"jwt_secret\")")
	fmt.Println("newSecretsManager.LoadKubernetesSecret(\"auth.api_key\", \"/etc/secrets\")")
} 