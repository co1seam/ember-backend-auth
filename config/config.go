package config

type App struct {
	Name    string `mapstructure:"APP_NAME"`
	Prefork bool   `mapstructure:"APP_PREFORK"`
}

type Database struct {
	Host string `mapstructure:"POSTGRES_HOST"`
	Port string `mapstructure:"POSTGRES_PORT"`
	User string `mapstructure:"POSTGRES_USER"`
	Pass string `mapstructure:"POSTGRES_PASS"`
	Name string `mapstructure:"POSTGRES_NAME"`
}

type SMTP struct {
	Host string `mapstructure:"SMTP_HOST"`
	Port string `mapstructure:"SMTP_PORT"`
	From string `mapstructure:"SMTP_FROM"`
}

type Token struct {
	Secret          string `mapstructure:"TOKEN_SECRET"`
	RefreshTokenTTL string `mapstructure:"TOKEN_REFRESH_TTL"`
	AccessTokenTTL  string `mapstructure:"TOKEN_ACCESS_TTL"`
}

type Redis struct {
	Host string `mapstructure:"REDIS_HOST"`
	Port string `mapstructure:"REDIS_PORT"`
}

type Config struct {
	App      App      `mapstructure:",squash"`
	Database Database `mapstructure:",squash"`
	SMTP     SMTP     `mapstructure:",squash"`
	Token    Token    `mapstructure:",squash"`
	Redis    Redis    `mapstructure:",squash"`
}
