package config_lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// LoadFromFile загружает конфигурацию из указанного файла
func LoadFromFile(filePath string) (*Config, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext == "" {
		return nil, fmt.Errorf("file extension required to determine config type")
	}
	
	configType := ext[1:] // Убираем точку из расширения
	configName := strings.TrimSuffix(filepath.Base(filePath), ext)
	configPath := filepath.Dir(filePath)
	
	options := Options{
		ConfigName: configName,
		ConfigPath: configPath,
		ConfigType: configType,
	}
	
	config := New(options)
	if err := config.Load(); err != nil {
		return nil, err
	}
	
	return config, nil
}

// LoadFromEnv загружает конфигурацию только из переменных окружения
func LoadFromEnv(prefix string) *Config {
	options := Options{
		EnvPrefix:    prefix,
		AutomaticEnv: true,
	}
	
	config := New(options)
	config.Load() // Здесь игнорируем ошибку, т.к. отсутствие файла конфига - не ошибка
	
	return config
}

// LoadFromJSON загружает конфигурацию из JSON строки
func LoadFromJSON(jsonStr string) (*Config, error) {
	options := Options{}
	config := New(options)
	
	if err := config.Load(); err != nil {
		return nil, err
	}
	
	if err := config.viper.ReadConfig(strings.NewReader(jsonStr)); err != nil {
		return nil, fmt.Errorf("error parsing JSON config: %w", err)
	}
	
	return config, nil
}

// LoadFromYAML загружает конфигурацию из YAML строки
func LoadFromYAML(yamlStr string) (*Config, error) {
	options := Options{
		ConfigType: "yaml",
	}
	config := New(options)
	
	if err := config.Load(); err != nil {
		return nil, err
	}
	
	if err := config.viper.ReadConfig(strings.NewReader(yamlStr)); err != nil {
		return nil, fmt.Errorf("error parsing YAML config: %w", err)
	}
	
	return config, nil
}

// LoadFromDotEnv загружает конфигурацию из .env файла
func LoadFromDotEnv(filePath string) (*Config, error) {
	if filePath == "" {
		filePath = ".env"
	}
	
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading .env file: %w", err)
	}
	
	options := Options{}
	config := New(options)
	
	// Обрабатываем .env файл строка за строкой
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// Пропускаем пустые строки и комментарии
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		// Разбиваем на ключ и значение
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		
		// Удаляем кавычки
		value = strings.Trim(value, `"'`)
		
		// Устанавливаем значение
		config.Set(key, value)
	}
	
	return config, nil
}

// SaveToFile сохраняет текущую конфигурацию в файл
func (c *Config) SaveToFile(filePath string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	// Определяем формат на основе расширения
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext == "" {
		return fmt.Errorf("file extension required to determine config type")
	}
	
	// Получаем все настройки
	settings := c.viper.AllSettings()
	
	var data []byte
	var err error
	
	// Сериализуем в зависимости от формата
	switch ext[1:] {
	case "json":
		data, err = json.MarshalIndent(settings, "", "  ")
	case "yaml", "yml":
		data, err = c.viper.WriteConfigAs(filePath)
		return err // WriteConfigAs возвращает nil, если запись прошла успешно
	case "toml":
		data, err = c.viper.WriteConfigAs(filePath)
		return err
	default:
		return fmt.Errorf("unsupported config format: %s", ext[1:])
	}
	
	if err != nil {
		return fmt.Errorf("error marshaling config: %w", err)
	}
	
	// Создаем директорию, если она не существует
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}
	
	// Записываем файл
	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}
	
	return nil
}

// LoadAndMerge загружает конфигурацию из нескольких источников с порядком приоритета
func LoadAndMerge(appName string, configPaths []string, envPrefix string) (*Config, error) {
	// Создаем основной конфиг
	options := Options{
		AppName:      appName,
		ConfigName:   "config",
		AutomaticEnv: true,
		EnvPrefix:    envPrefix,
	}
	
	if envPrefix == "" && appName != "" {
		options.EnvPrefix = strings.ToUpper(appName)
	}
	
	config := New(options)
	
	// Добавляем пути поиска конфигурации
	for _, path := range configPaths {
		config.viper.AddConfigPath(path)
	}
	
	// Пытаемся загрузить конфигурацию
	if err := config.Load(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Если ошибка не связана с отсутствием файла, возвращаем её
			return nil, err
		}
	}
	
	return config, nil
}

// MergeConfig объединяет два экземпляра Config
func (c *Config) MergeConfig(other *Config) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	// Получаем все настройки из другого конфига
	settings := other.viper.AllSettings()
	
	// Объединяем настройки
	for key, value := range settings {
		c.viper.Set(key, value)
	}
} 