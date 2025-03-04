package repo_lib

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Model базовая модель для всех сущностей
type Model struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// Repository интерфейс базового репозитория
type Repository[T any] interface {
	// Create создает новую запись
	Create(ctx context.Context, entity *T) error
	
	// CreateBatch создает множество записей
	CreateBatch(ctx context.Context, entities []*T) error
	
	// Update обновляет запись
	Update(ctx context.Context, entity *T) error
	
	// Delete удаляет запись
	Delete(ctx context.Context, id uint) error
	
	// DeleteBatch удаляет множество записей
	DeleteBatch(ctx context.Context, ids []uint) error
	
	// FindByID находит запись по ID
	FindByID(ctx context.Context, id uint) (*T, error)
	
	// FindAll находит все записи с пагинацией
	FindAll(ctx context.Context, page, size int) ([]*T, int64, error)
	
	// FindByCondition находит записи по условию
	FindByCondition(ctx context.Context, condition map[string]interface{}) ([]*T, error)
	
	// Count возвращает количество записей
	Count(ctx context.Context) (int64, error)
	
	// Transaction выполняет операции в транзакции
	Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error
}

// BaseRepository базовая реализация репозитория
type BaseRepository[T any] struct {
	db *gorm.DB
}

// NewBaseRepository создает новый экземпляр базового репозитория
func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{
		db: db,
	}
}

// Create создает новую запись
func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

// CreateBatch создает множество записей
func (r *BaseRepository[T]) CreateBatch(ctx context.Context, entities []*T) error {
	return r.db.WithContext(ctx).CreateInBatches(entities, 100).Error
}

// Update обновляет запись
func (r *BaseRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

// Delete удаляет запись
func (r *BaseRepository[T]) Delete(ctx context.Context, id uint) error {
	var entity T
	result := r.db.WithContext(ctx).Delete(&entity, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("record not found")
	}
	return nil
}

// DeleteBatch удаляет множество записей
func (r *BaseRepository[T]) DeleteBatch(ctx context.Context, ids []uint) error {
	var entity T
	return r.db.WithContext(ctx).Delete(&entity, ids).Error
}

// FindByID находит запись по ID
func (r *BaseRepository[T]) FindByID(ctx context.Context, id uint) (*T, error) {
	var entity T
	if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("record with id %d not found", id)
		}
		return nil, err
	}
	return &entity, nil
}

// FindAll находит все записи с пагинацией
func (r *BaseRepository[T]) FindAll(ctx context.Context, page, size int) ([]*T, int64, error) {
	var entities []*T
	var total int64

	// Получаем общее количество записей
	if err := r.db.WithContext(ctx).Model(new(T)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Получаем записи с пагинацией
	offset := (page - 1) * size
	if err := r.db.WithContext(ctx).Offset(offset).Limit(size).Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

// FindByCondition находит записи по условию
func (r *BaseRepository[T]) FindByCondition(ctx context.Context, condition map[string]interface{}) ([]*T, error) {
	var entities []*T
	if err := r.db.WithContext(ctx).Where(condition).Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

// Count возвращает количество записей
func (r *BaseRepository[T]) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(new(T)).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Transaction выполняет операции в транзакции
func (r *BaseRepository[T]) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fn)
} 