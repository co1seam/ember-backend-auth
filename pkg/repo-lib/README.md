# Repo-Lib

Гибкая библиотека для работы с базами данных в Go, основанная на [GORM](https://gorm.io) с поддержкой кэширования через [Redis](https://redis.io).

## Особенности

- **Универсальный интерфейс** - единый интерфейс для работы с разными базами данных
- **Поддержка нескольких СУБД** - PostgreSQL, MySQL, SQLite
- **Типобезопасность** - использование дженериков для типобезопасной работы с моделями
- **Транзакции** - поддержка транзакций
- **Пагинация** - встроенная поддержка пагинации
- **Гибкая конфигурация** - настройка подключения и пула соединений
- **Логирование** - настраиваемое логирование запросов
- **Кэширование** - поддержка кэширования через Redis
- **Автоматическая инвалидация** - автоматическое обновление кэша при изменениях

## Установка

```bash
go get github.com/co1seam/tuneflow-backend-auth/pkg/repo-lib
```

## Быстрый старт

```go
package main

import (
    "context"
    "log"
    
    repo_lib "github.com/co1seam/tuneflow-backend-auth/pkg/repo-lib"
)

// User модель пользователя
type User struct {
    repo_lib.Model
    Name  string
    Email string
}

func main() {
    // Конфигурация базы данных
    dbConfig := repo_lib.DefaultConfig()
    dbConfig.Database = "example"
    
    // Подключаемся к базе данных
    db, err := repo_lib.Connect(dbConfig)
    if err != nil {
        log.Fatal(err)
    }
    defer repo_lib.Disconnect(db)
    
    // Конфигурация Redis
    cacheConfig := repo_lib.DefaultCacheConfig()
    
    // Подключаемся к Redis
    cache, err := repo_lib.NewRedisCache(cacheConfig)
    if err != nil {
        log.Fatal(err)
    }
    defer cache.Close()
    
    // Создаем базовый репозиторий
    baseRepo := repo_lib.NewBaseRepository[User](db)
    
    // Создаем кэшированный репозиторий
    userRepo := repo_lib.NewCachedRepository[User](baseRepo, cache)
    
    // Используем репозиторий
    ctx := context.Background()
    user := &User{
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    if err := userRepo.Create(ctx, user); err != nil {
        log.Fatal(err)
    }
}
```

## Конфигурация

### Настройки базы данных

```go
cfg := repo_lib.DefaultConfig()
```

Настройки по умолчанию:
- Type: PostgreSQL
- Host: localhost
- Port: 5432
- User: postgres
- Password: postgres
- Database: postgres
- SSLMode: disable
- MaxIdleConns: 10
- MaxOpenConns: 100
- ConnMaxLifetime: 1 час
- LogLevel: Info
- SlowThreshold: 1 секунда

### Настройки Redis

```go
cfg := repo_lib.DefaultCacheConfig()
```

Настройки по умолчанию:
- Host: localhost
- Port: 6379
- Password: ""
- DB: 0
- MaxRetries: 3
- DialTimeout: 5 секунд
- ReadTimeout: 3 секунды
- WriteTimeout: 3 секунды
- PoolSize: 10
- MinIdleConns: 5
- TTL: 1 час

## Использование

### Создание репозитория

```go
// Базовый репозиторий
baseRepo := repo_lib.NewBaseRepository[User](db)

// Кэшированный репозиторий
cachedRepo := repo_lib.NewCachedRepository[User](baseRepo, cache)

// Пользовательский репозиторий
type UserRepository struct {
    *repo_lib.CachedRepository[User]
}

func NewUserRepository(db *gorm.DB, cache repo_lib.Cache) *UserRepository {
    baseRepo := repo_lib.NewBaseRepository[User](db)
    cachedRepo := repo_lib.NewCachedRepository[User](baseRepo, cache)
    return &UserRepository{
        CachedRepository: cachedRepo,
    }
}
```

### Базовые операции

```go
// Создание (сохраняется в БД и кэш)
user := &User{Name: "John"}
err := repo.Create(ctx, user)

// Обновление (обновляется в БД и кэш)
user.Name = "John Smith"
err := repo.Update(ctx, user)

// Удаление (удаляется из БД и кэша)
err := repo.Delete(ctx, user.ID)

// Поиск по ID (сначала ищет в кэше, потом в БД)
user, err := repo.FindByID(ctx, 1)
```

### Работа с кэшем

```go
// Прямая работа с кэшем
err := cache.Set(ctx, "key", value, time.Hour)
err := cache.Get(ctx, "key", &value)
err := cache.Delete(ctx, "key")
err := cache.Clear(ctx)
```

### Транзакции

```go
err := repo.Transaction(ctx, func(tx *gorm.DB) error {
    // Операции в транзакции
    user := &User{Name: "John"}
    if err := tx.Create(user).Error; err != nil {
        return err
    }
    
    return nil
})
```

## Примеры

Полные примеры использования можно найти в директории `example/`:

- `main.go` - базовый пример использования
- `cache_example.go` - пример использования кэширования
- `custom_repo.go` - пример пользовательского репозитория
- `transactions.go` - пример работы с транзакциями
- `pagination.go` - пример пагинации

## Лицензия

MIT 