# Config-Lib

Гибкая библиотека для работы с конфигурацией в Go, следующая принципам 12-factor приложений.

## Особенности

- **Следование принципам 12-factor app** - конфигурация строго отделена от кода
- **Множество форматов конфигурации** - поддержка YAML, JSON, TOML и других форматов
- **Переменные окружения** - легкое переопределение через переменные окружения
- **Многоуровневая конфигурация** - поддержка разных конфигураций для разных окружений
- **Автообновление конфигурации** - наблюдение за изменениями файлов конфигурации
- **Валидация конфигурации** - проверка наличия обязательных параметров
- **Удобный интерфейс** - множество вспомогательных методов
- **Безопасность** - поддержка загрузки секретов из файлов
- **Конкурентная безопасность** - защита от гонок данных при доступе из множества горутин

## Установка

```bash
go get github.com/co1seam/tuneflow-backend-auth/pkg/config-lib
```

## Быстрый старт

```go
package main

import (
    "fmt"
    "log"
    
    config_lib "github.com/co1seam/tuneflow-backend-auth/pkg/config-lib"
)

func main() {
    // Создаем конфигурацию, соответствующую 12-factor app
    config := config_lib.TwelveFactorConfig("myapp")
    
    // Проверяем окружение
    if config.IsProduction() {
        fmt.Println("Running in production mode")
    } else {
        fmt.Println("Running in development mode")
    }
    
    // Получаем значения
    port := config.GetString("app.port")
    host := config.GetString("app.host")
    
    fmt.Printf("Server will run on %s:%s\n", host, port)
}
```

## Загрузка конфигурации из разных источников

### Из файла

```go
config, err := config_lib.LoadFromFile("config.yaml")
if err != nil {
    log.Fatalf("Failed to load config: %v", err)
}
```

### Из переменных окружения

```go
config := config_lib.LoadFromEnv("MYAPP")
```

### Из JSON или YAML строки

```go
jsonConfig := `{"server": {"port": 8080}}`
config, err := config_lib.LoadFromJSON(jsonConfig)
if err != nil {
    log.Fatalf("Failed to load config: %v", err)
}
```

### Из .env файла

```go
config, err := config_lib.LoadFromDotEnv(".env")
if err != nil {
    log.Fatalf("Failed to load config: %v", err)
}
```

## Поддержка 12-factor приложений

Config-Lib следует принципам [12-factor app](https://12factor.net/), особенно в отношении [конфигурации](https://12factor.net/config):

1. **Строгое разделение конфигурации и кода**:
   ```go
   // Конфигурация хранится отдельно от кода
   config := config_lib.DefaultConfig("myapp")
   ```

2. **Использование переменных окружения**:
   ```go
   // Автоматическое чтение из переменных окружения
   options := config_lib.Options{
       AutomaticEnv: true,
       EnvPrefix: "MYAPP",
   }
   ```

3. **Конфигурация для разных окружений**:
   ```go
   // Загрузка конфигурации в зависимости от окружения
   config, err := config_lib.LoadEnvBasedConfig("myapp")
   ```

## Отображение конфигурации на структуры

```go
type DatabaseConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    Username string `mapstructure:"username"`
    Password string `mapstructure:"password"`
}

var dbConfig DatabaseConfig
if err := config.UnmarshalKey("database", &dbConfig); err != nil {
    log.Fatalf("Failed to unmarshal db config: %v", err)
}
```

## Работа с секретами

```go
// Загрузка секрета из переменной окружения или файла
// Теперь возвращает ошибку, если секрет не найден
err := config.LoadSecret("database.password", "DB_PASSWORD", "/run/secrets/db_password")
if err != nil {
    log.Printf("Couldn't load secret: %v", err)
}
```

## Динамическое обновление конфигурации

```go
config := config_lib.New(config_lib.Options{
    AppName:     "myapp",
    WatchConfig: true,
})

config.OnChange(func() {
    fmt.Println("Config has been updated!")
    // Получаем обновленные значения
    newPort := config.GetInt("server.port")
    fmt.Printf("New port: %d\n", newPort)
})
```

## Пример конфигурационного файла

### config.yaml

```yaml
app:
  name: "MyApp"
  env: "development"
  log_level: "debug"

server:
  host: "localhost"
  port: 8080
  read_timeout: "5s"
  write_timeout: "10s"
  debug: true

database:
  host: "localhost"
  port: 5432
  username: "postgres"
  password: "password"
  name: "myapp"
  ssl: false
  timeout: 10

features:
  - "auth"
  - "api"
  - "admin"
```

## Полный пример использования

```go
package main

import (
    "fmt"
    "log"
    "time"
    
    config_lib "github.com/co1seam/tuneflow-backend-auth/pkg/config-lib"
)

// Структура конфигурации приложения
type AppConfig struct {
    Name        string        `mapstructure:"name"`
    Environment string        `mapstructure:"env"`
    Server      ServerConfig  `mapstructure:"server"`
    Database    DatabaseConfig `mapstructure:"database"`
    Features    []string      `mapstructure:"features"`
}

type ServerConfig struct {
    Host         string        `mapstructure:"host"`
    Port         int           `mapstructure:"port"`
    ReadTimeout  time.Duration `mapstructure:"read_timeout"`
    WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    Username string `mapstructure:"username"`
    Password string `mapstructure:"password"`
    Name     string `mapstructure:"name"`
}

func main() {
    // Создаем конфигурацию, соответствующую 12-factor app
    config, err := config_lib.LoadEnvBasedConfig("myapp")
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    // Устанавливаем значения по умолчанию
    config.SetDefault("server.port", 8080)
    config.SetDefault("server.host", "0.0.0.0")
    
    // Загружаем секреты
    err = config.LoadSecret("database.password", "DB_PASSWORD", "/run/secrets/db_password")
    if err != nil {
        log.Printf("Warning: %v", err)
    }
    
    // Отображаем конфигурацию на структуру
    var appConfig AppConfig
    if err := config.Unmarshal(&appConfig); err != nil {
        log.Fatalf("Failed to unmarshal config: %v", err)
    }
    
    // Выводим конфигурацию
    fmt.Printf("App: %s (env: %s)\n", appConfig.Name, appConfig.Environment)
    fmt.Printf("Server: %s:%d\n", appConfig.Server.Host, appConfig.Server.Port)
    fmt.Printf("Database: %s@%s:%d/%s\n", 
        appConfig.Database.Username, 
        appConfig.Database.Host, 
        appConfig.Database.Port, 
        appConfig.Database.Name,
    )
    fmt.Printf("Features: %v\n", appConfig.Features)
    
    // Подписываемся на изменения конфигурации
    config.OnChange(func() {
        // Обновляем конфигурацию при изменении
        config.Unmarshal(&appConfig)
        fmt.Println("Config has been updated!")
    })
    
    // Запускаем сервер...
}
```

## Создание примеров конфигурационных файлов

Библиотека позволяет сохранить текущую конфигурацию в файл, что удобно для создания шаблонов:

```go
config := config_lib.DefaultConfig("myapp")
config.Set("app.name", "MyApp")
config.Set("app.env", "development")
config.Set("server.port", 8080)

// Сохраняем конфигурацию в файл
if err := config.SaveToFile("config.example.yaml"); err != nil {
    log.Fatalf("Failed to save config: %v", err)
}
```

## Лицензия

MIT 