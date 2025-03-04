package repo_lib

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
)

// CacheConfig конфигурация Redis
type CacheConfig struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	Password     string        `yaml:"password"`
	DB           int          `yaml:"db"`
	MaxRetries   int          `yaml:"max_retries"`
	DialTimeout  time.Duration `yaml:"dial_timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	PoolSize     int          `yaml:"pool_size"`
	MinIdleConns int          `yaml:"min_idle_conns"`
	TTL          time.Duration `yaml:"ttl"`
}

// DefaultCacheConfig возвращает конфигурацию Redis по умолчанию
func DefaultCacheConfig() *CacheConfig {
	return &CacheConfig{
		Host:         "localhost",
		Port:         6379,
		Password:     "",
		DB:           0,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 5,
		TTL:          time.Hour,
	}
}

// Cache интерфейс для работы с кэшем
type Cache interface {
	// Get получает значение из кэша
	Get(ctx context.Context, key string, value interface{}) error
	
	// Set устанавливает значение в кэш
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	
	// Delete удаляет значение из кэша
	Delete(ctx context.Context, key string) error
	
	// Clear очищает кэш
	Clear(ctx context.Context) error
}

// RedisCache реализация кэша на Redis
type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

// NewRedisCache создает новый экземпляр Redis кэша
func NewRedisCache(cfg *CacheConfig) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	})

	// Проверяем подключение
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisCache{
		client: client,
		ttl:    cfg.TTL,
	}, nil
}

// Get получает значение из кэша
func (c *RedisCache) Get(ctx context.Context, key string, value interface{}) error {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("key %s not found in cache", key)
		}
		return fmt.Errorf("failed to get value from cache: %w", err)
	}

	if err := msgpack.Unmarshal(data, value); err != nil {
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}

	return nil
}

// Set устанавливает значение в кэш
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := msgpack.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	if ttl == 0 {
		ttl = c.ttl
	}

	if err := c.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set value in cache: %w", err)
	}

	return nil
}

// Delete удаляет значение из кэша
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	if err := c.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete value from cache: %w", err)
	}
	return nil
}

// Clear очищает кэш
func (c *RedisCache) Clear(ctx context.Context) error {
	if err := c.client.FlushDB(ctx).Err(); err != nil {
		return fmt.Errorf("failed to clear cache: %w", err)
	}
	return nil
}

// Close закрывает соединение с Redis
func (c *RedisCache) Close() error {
	return c.client.Close()
}

// CachedRepository обертка над репозиторием с поддержкой кэширования
type CachedRepository[T any] struct {
	repo  Repository[T]
	cache Cache
}

// NewCachedRepository создает новый репозиторий с кэшированием
func NewCachedRepository[T any](repo Repository[T], cache Cache) *CachedRepository[T] {
	return &CachedRepository[T]{
		repo:  repo,
		cache: cache,
	}
}

// getCacheKey генерирует ключ для кэша
func getCacheKey(prefix string, id interface{}) string {
	return fmt.Sprintf("%s:%v", prefix, id)
}

// Create создает новую запись
func (r *CachedRepository[T]) Create(ctx context.Context, entity *T) error {
	if err := r.repo.Create(ctx, entity); err != nil {
		return err
	}

	// Кэшируем созданную запись
	key := getCacheKey("entity", getEntityID(entity))
	return r.cache.Set(ctx, key, entity, 0)
}

// Update обновляет запись
func (r *CachedRepository[T]) Update(ctx context.Context, entity *T) error {
	if err := r.repo.Update(ctx, entity); err != nil {
		return err
	}

	// Обновляем кэш
	key := getCacheKey("entity", getEntityID(entity))
	return r.cache.Set(ctx, key, entity, 0)
}

// Delete удаляет запись
func (r *CachedRepository[T]) Delete(ctx context.Context, id uint) error {
	if err := r.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Удаляем из кэша
	key := getCacheKey("entity", id)
	return r.cache.Delete(ctx, key)
}

// FindByID находит запись по ID
func (r *CachedRepository[T]) FindByID(ctx context.Context, id uint) (*T, error) {
	var entity T
	key := getCacheKey("entity", id)

	// Пытаемся получить из кэша
	err := r.cache.Get(ctx, key, &entity)
	if err == nil {
		return &entity, nil
	}

	// Если нет в кэше, получаем из БД
	result, err := r.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Кэшируем результат
	if err := r.cache.Set(ctx, key, result, 0); err != nil {
		return nil, err
	}

	return result, nil
}

// getEntityID получает ID сущности через рефлексию
func getEntityID(entity interface{}) uint {
	v := reflect.ValueOf(entity)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return 0
	}
	if f := v.FieldByName("ID"); f.IsValid() {
		return uint(f.Uint())
	}
	return 0
} 