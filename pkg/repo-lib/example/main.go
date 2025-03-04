package main

import (
	"context"
	"fmt"
	"log"

	repo_lib "github.com/co1seam/tuneflow-backend-auth/pkg/repo-lib"
)

// User модель пользователя
type User struct {
	repo_lib.Model
	Name     string `gorm:"size:255;not null" json:"name"`
	Email    string `gorm:"size:255;not null;unique" json:"email"`
	Password string `gorm:"size:255;not null" json:"-"`
}

// UserRepository репозиторий для работы с пользователями
type UserRepository struct {
	*repo_lib.BaseRepository[User]
}

// NewUserRepository создает новый репозиторий пользователей
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		BaseRepository: repo_lib.NewBaseRepository[User](db),
	}
}

func main() {
	// Создаем конфигурацию базы данных
	cfg := repo_lib.DefaultConfig()
	cfg.Database = "example"

	// Подключаемся к базе данных
	db, err := repo_lib.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer repo_lib.Disconnect(db)

	// Автомиграция
	if err := db.AutoMigrate(&User{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Создаем репозиторий
	userRepo := NewUserRepository(db)

	// Создаем контекст
	ctx := context.Background()

	// Создаем пользователя
	user := &User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "secret",
	}

	if err := userRepo.Create(ctx, user); err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	fmt.Printf("Created user: %+v\n", user)

	// Находим пользователя по ID
	found, err := userRepo.FindByID(ctx, user.ID)
	if err != nil {
		log.Fatalf("Failed to find user: %v", err)
	}
	fmt.Printf("Found user: %+v\n", found)

	// Обновляем пользователя
	user.Name = "John Smith"
	if err := userRepo.Update(ctx, user); err != nil {
		log.Fatalf("Failed to update user: %v", err)
	}
	fmt.Printf("Updated user: %+v\n", user)

	// Получаем всех пользователей
	users, total, err := userRepo.FindAll(ctx, 1, 10)
	if err != nil {
		log.Fatalf("Failed to find users: %v", err)
	}
	fmt.Printf("Found %d users, total: %d\n", len(users), total)

	// Находим пользователей по условию
	condition := map[string]interface{}{
		"name": "John Smith",
	}
	filtered, err := userRepo.FindByCondition(ctx, condition)
	if err != nil {
		log.Fatalf("Failed to find users by condition: %v", err)
	}
	fmt.Printf("Found users by condition: %+v\n", filtered)

	// Удаляем пользователя
	if err := userRepo.Delete(ctx, user.ID); err != nil {
		log.Fatalf("Failed to delete user: %v", err)
	}
	fmt.Printf("Deleted user with ID: %d\n", user.ID)

	// Пример использования транзакции
	err = userRepo.Transaction(ctx, func(tx *gorm.DB) error {
		// Создаем нового пользователя в транзакции
		newUser := &User{
			Name:     "Jane Doe",
			Email:    "jane@example.com",
			Password: "secret",
		}
		
		if err := tx.Create(newUser).Error; err != nil {
			return err
		}
		
		// Обновляем другого пользователя
		if err := tx.Model(&User{}).Where("email = ?", "other@example.com").
			Update("name", "Other Name").Error; err != nil {
			return err
		}
		
		return nil
	})
	if err != nil {
		log.Fatalf("Transaction failed: %v", err)
	}
} 