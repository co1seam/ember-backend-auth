package config_lib

// Экспортируемые типы и константы
// Этот файл содержит часто используемые типы и константы,
// которые удобны для быстрого импорта пользователями

// Константы для уровней окружения
const (
	// Окружение разработки
	Development = EnvDevelopment
	
	// Тестовое окружение
	Test = EnvTest
	
	// Предпродакшн окружение
	Staging = EnvStaging
	
	// Продакшн окружение
	Production = EnvProduction
)

// Константы для уровней логирования
const (
	// Уровень логирования Debug
	LogDebug = LogLevelDebug
	
	// Уровень логирования Info
	LogInfo = LogLevelInfo
	
	// Уровень логирования Warning
	LogWarn = LogLevelWarn
	
	// Уровень логирования Error
	LogError = LogLevelError
)

// Константы для переменных окружения
const (
	// Переменная окружения для среды выполнения
	EnvVar = AppEnv
	
	// Переменная окружения для порта
	PortVar = AppPort
	
	// Переменная окружения для хоста
	HostVar = AppHost
	
	// Переменная окружения для режима отладки
	DebugVar = AppDebug
)

// Option представляет функцию-опцию для настройки параметров конфигурации
type Option func(*Options)

// WithAppName устанавливает имя приложения
func WithAppName(appName string) Option {
	return func(options *Options) {
		options.AppName = appName
	}
}

// WithConfigPath устанавливает путь к конфигурационным файлам
func WithConfigPath(configPath string) Option {
	return func(options *Options) {
		options.ConfigPath = configPath
	}
}

// WithConfigName устанавливает имя конфигурационного файла
func WithConfigName(configName string) Option {
	return func(options *Options) {
		options.ConfigName = configName
	}
}

// WithConfigType устанавливает тип конфигурационного файла
func WithConfigType(configType string) Option {
	return func(options *Options) {
		options.ConfigType = configType
	}
}

// WithEnvPrefix устанавливает префикс для переменных окружения
func WithEnvPrefix(envPrefix string) Option {
	return func(options *Options) {
		options.EnvPrefix = envPrefix
	}
}

// WithAutomaticEnv включает автоматическое чтение переменных окружения
func WithAutomaticEnv(enable bool) Option {
	return func(options *Options) {
		options.AutomaticEnv = enable
	}
}

// WithWatchConfig включает наблюдение за изменениями конфигурационного файла
func WithWatchConfig(enable bool) Option {
	return func(options *Options) {
		options.WatchConfig = enable
	}
}

// WithDefaults устанавливает значения по умолчанию
func WithDefaults(defaults map[string]interface{}) Option {
	return func(options *Options) {
		options.Defaults = defaults
	}
}

// NewConfig создает новый экземпляр Config с применением опций
func NewConfig(opts ...Option) *Config {
	options := Options{
		AutomaticEnv: true,
		KeyDelimiter: ".",
	}
	
	for _, opt := range opts {
		opt(&options)
	}
	
	return New(options)
} 