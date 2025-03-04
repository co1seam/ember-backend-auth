package main

import (
	"context"
	"fmt"
	"log"
	"time"

	repo_lib "github.com/co1seam/tuneflow-backend-auth/pkg/repo-lib"
)

func main() {
	// Создаем конфигурацию базы данных
	dbConfig := repo_lib.DefaultConfig()
	dbConfig.Database = "example"

	// Подключаемся к базе данных
	db, err := repo_lib.Connect(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer repo_lib.Disconnect(db)

	// Создаем конфигурацию Redis
	cacheConfig := repo_lib.DefaultCacheConfig()
	cacheConfig.TTL = 5 * time.Minute

	// Подключаемся к Redis
	cache, err := repo_lib.NewRedisCache(cacheConfig)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer cache.Close()

	// Создаем базовый репозиторий
	baseRepo := repo_lib.NewBaseRepository[User](db)

	// Создаем кэшированный репозиторий
	userRepo := repo_lib.NewCachedRepository[User](baseRepo, cache)

	// Создаем контекст
	ctx := context.Background()

	// Создаем пользователя
	user := &User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "secret",
	}

	// Создаем запись (сохранится в БД и в кэш)
	if err := userRepo.Create(ctx, user); err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	fmt.Printf("Created user: %+v\n", user)

	// Первый поиск (загрузится из БД и сохранится в кэш)
	found1, err := userRepo.FindByID(ctx, user.ID)
	if err != nil {
		log.Fatalf("Failed to find user: %v", err)
	}
	fmt.Printf("Found user (from DB): %+v\n", found1)

	// Второй поиск (загрузится из кэша)
	found2, err := userRepo.FindByID(ctx, user.ID)
	if err != nil {
		log.Fatalf("Failed to find user: %v", err)
	}
	fmt.Printf("Found user (from cache): %+v\n", found2)

	// Обновляем пользователя (обновится в БД и в кэше)
	user.Name = "John Smith"
	if err := userRepo.Update(ctx, user); err != nil {
		log.Fatalf("Failed to update user: %v", err)
	}
	fmt.Printf("Updated user: %+v\n", user)

	// Поиск после обновления (загрузится из кэша)
	found3, err := userRepo.FindByID(ctx, user.ID)
	if err != nil {
		log.Fatalf("Failed to find user: %v", err)
	}
	fmt.Printf("Found user after update (from cache): %+v\n", found3)

	// Удаляем пользователя (удалится из БД и из кэша)
	if err := userRepo.Delete(ctx, user.ID); err != nil {
		log.Fatalf("Failed to delete user: %v", err)
	}
	fmt.Printf("Deleted user with ID: %d\n", user.ID)

	// Пытаемся найти удаленного пользователя
	found4, err := userRepo.FindByID(ctx, user.ID)
	if err != nil {
		fmt.Printf("Expected error: %v\n", err)
	} else {
		fmt.Printf("User should not be found: %+v\n", found4)
	}
} 