package config_lib

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// SecretsManager управляет безопасными секретами
type SecretsManager struct {
	config *Config
}

// NewSecretsManager создает новый менеджер секретов
func NewSecretsManager(config *Config) *SecretsManager {
	return &SecretsManager{
		config: config,
	}
}

// LoadSecretFromEnv загружает секрет из переменной окружения
func (sm *SecretsManager) LoadSecretFromEnv(key, envVar string) bool {
	if sm.config.IsSet(key) {
		return true
	}
	
	if value := os.Getenv(envVar); value != "" {
		sm.config.Set(key, value)
		return true
	}
	
	return false
}

// LoadSecretFromFile загружает секрет из файла
func (sm *SecretsManager) LoadSecretFromFile(key, filePath string) error {
	if sm.config.IsSet(key) {
		return nil
	}
	
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading secret file '%s': %w", filePath, err)
	}
	
	sm.config.Set(key, strings.TrimSpace(string(data)))
	return nil
}

// LoadDockerSecret загружает секрет из Docker Swarm секрета
func (sm *SecretsManager) LoadDockerSecret(key, secretName string) error {
	if sm.config.IsSet(key) {
		return nil
	}
	
	// Docker Swarm монтирует секреты в /run/secrets/<secret_name>
	secretPath := filepath.Join("/run/secrets", secretName)
	return sm.LoadSecretFromFile(key, secretPath)
}

// LoadKubernetesSecret загружает секрет из Kubernetes секрета
func (sm *SecretsManager) LoadKubernetesSecret(key, secretPath string) error {
	if sm.config.IsSet(key) {
		return nil
	}
	
	// В Kubernetes секреты обычно монтируются как файлы в указанную директорию
	if secretPath == "" {
		secretPath = "/etc/secrets"
	}
	
	// Имя файла совпадает с ключом
	secretFile := filepath.Join(secretPath, filepath.Base(key))
	return sm.LoadSecretFromFile(key, secretFile)
}

// Encrypt шифрует строку с помощью AES-GCM
func (sm *SecretsManager) Encrypt(plaintext string, passphrase string) (string, error) {
	// Создаем шифр
	block, err := aes.NewCipher([]byte(createHash(passphrase)))
	if err != nil {
		return "", err
	}
	
	// Создаем GCM (Galois/Counter Mode)
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	
	// Создаем nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	
	// Шифруем
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	
	// Кодируем в base64
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt расшифровывает строку с помощью AES-GCM
func (sm *SecretsManager) Decrypt(cryptoText string, passphrase string) (string, error) {
	// Декодируем base64
	ciphertext, err := base64.StdEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}
	
	// Создаем шифр
	block, err := aes.NewCipher([]byte(createHash(passphrase)))
	if err != nil {
		return "", err
	}
	
	// Создаем GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	
	// Проверяем длину
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	
	// Извлекаем nonce и расшифровываем
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	
	return string(plaintext), nil
}

// CreateEncryptedConfig создает зашифрованную конфигурацию
func (sm *SecretsManager) CreateEncryptedConfig(filePath, passphrase string) error {
	// Получаем настройки
	settings := sm.config.viper.AllSettings()
	
	// Создаем временный конфиг для сохранения зашифрованных данных
	encConfig := DefaultConfig("encrypted")
	
	// Шифруем каждый секретный ключ
	for key, value := range flattenMap(settings, "") {
		// Определяем, является ли ключ секретным (содержит слова password, secret, key и т.д.)
		if isSecretKey(key) {
			if strValue, ok := value.(string); ok {
				encValue, err := sm.Encrypt(strValue, passphrase)
				if err != nil {
					return fmt.Errorf("error encrypting value for key '%s': %w", key, err)
				}
				encConfig.Set(key, fmt.Sprintf("ENC(%s)", encValue))
			}
		} else {
			encConfig.Set(key, value)
		}
	}
	
	// Сохраняем зашифрованный конфиг
	return encConfig.SaveToFile(filePath)
}

// LoadEncryptedConfig загружает зашифрованную конфигурацию
func (sm *SecretsManager) LoadEncryptedConfig(filePath, passphrase string) error {
	// Загружаем конфигурацию из файла
	encConfig, err := LoadFromFile(filePath)
	if err != nil {
		return fmt.Errorf("error loading encrypted config: %w", err)
	}
	
	// Получаем все настройки
	settings := encConfig.viper.AllSettings()
	
	// Расшифровываем каждый зашифрованный ключ
	for key, value := range flattenMap(settings, "") {
		if strValue, ok := value.(string); ok {
			// Проверяем, зашифровано ли значение
			if strings.HasPrefix(strValue, "ENC(") && strings.HasSuffix(strValue, ")") {
				// Извлекаем зашифрованное значение
				encValue := strValue[4 : len(strValue)-1]
				
				// Расшифровываем
				decValue, err := sm.Decrypt(encValue, passphrase)
				if err != nil {
					return fmt.Errorf("error decrypting value for key '%s': %w", key, err)
				}
				
				// Устанавливаем расшифрованное значение
				sm.config.Set(key, decValue)
			} else {
				// Устанавливаем незашифрованное значение
				sm.config.Set(key, strValue)
			}
		} else {
			// Устанавливаем нестроковое значение
			sm.config.Set(key, value)
		}
	}
	
	return nil
}

// Вспомогательные функции

// createHash создает хеш из пароля для использования в качестве ключа шифрования
func createHash(key string) []byte {
	// Простой способ получить 32 байта (для AES-256)
	// В реальных приложениях лучше использовать PBKDF2 или другие функции для получения ключа
	hash := make([]byte, 32)
	copy(hash, []byte(key))
	for i := len(key); i < 32; i++ {
		hash[i] = byte(i)
	}
	return hash
}

// isSecretKey проверяет, является ли ключ секретным
func isSecretKey(key string) bool {
	// Ключи, содержащие эти слова, считаются секретными
	secretWords := []string{
		"password", "passwd", "secret", "key", "token", "auth", "credential",
	}
	
	lowerKey := strings.ToLower(key)
	for _, word := range secretWords {
		if strings.Contains(lowerKey, word) {
			return true
		}
	}
	
	return false
}

// flattenMap преобразует вложенную карту в плоскую карту с составными ключами
func flattenMap(m map[string]interface{}, prefix string) map[string]interface{} {
	result := make(map[string]interface{})
	
	for k, v := range m {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}
		
		switch val := v.(type) {
		case map[string]interface{}:
			// Рекурсивно обрабатываем вложенные карты
			for k2, v2 := range flattenMap(val, key) {
				result[k2] = v2
			}
		case []interface{}:
			// Обрабатываем слайсы, но не шифруем их содержимое
			result[key] = val
		default:
			result[key] = val
		}
	}
	
	return result
} 