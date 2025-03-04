package repo_lib

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DBType тип базы данных
type DBType string

const (
	PostgreSQL DBType = "postgres"
	MySQL     DBType = "mysql"
	SQLite    DBType = "sqlite"
)

// Config конфигурация базы данных
type Config struct {
	Type     DBType `yaml:"type"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	SSLMode  string `yaml:"ssl_mode"`
	
	// Дополнительные настройки
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
	
	// Настройки логирования
	LogLevel  logger.LogLevel `yaml:"log_level"`
	SlowThreshold time.Duration `yaml:"slow_threshold"`
}

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() *Config {
	return &Config{
		Type:     PostgreSQL,
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		Database: "postgres",
		SSLMode:  "disable",
		
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: time.Hour,
		
		LogLevel:      logger.Info,
		SlowThreshold: time.Second,
	}
}

// Connect устанавливает соединение с базой данных
func Connect(cfg *Config) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch cfg.Type {
	case PostgreSQL:
		dsn := fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode,
		)
		dialector = postgres.Open(dsn)
	
	case MySQL:
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database,
		)
		dialector = mysql.Open(dsn)
	
	case SQLite:
		dialector = sqlite.Open(cfg.Database)
	
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Type)
	}

	// Настройки GORM
	config := &gorm.Config{
		Logger: logger.New(
			logger.Default.Writer(),
			logger.Config{
				SlowThreshold:             cfg.SlowThreshold,
				LogLevel:                  cfg.LogLevel,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
		),
	}

	// Устанавливаем соединение
	db, err := gorm.Open(dialector, config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Настраиваем пул соединений
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	return db, nil
}

// Disconnect закрывает соединение с базой данных
func Disconnect(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}
	
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}
	
	return nil
} 