package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"time"
)

const (
	defaultHTTPPort               = "8000"
	defaultHTTPRWTimeout          = 10 * time.Second
	defaultHTTPMaxHeaderMegabytes = 1
	defaultAccessTokenTTL         = 15 * time.Minute
	defaultRefreshTokenTTL        = 24 * time.Hour * 30
	defaultSecretCodeTTL          = 2 * time.Minute

	EnvLocal = "local"
	Prod     = "prod"
)

type (
	Config struct {
		Environment string
		Postgres    PostgresConfig
		HTTP        HTTPConfig
		Auth        AuthConfig
		Redis       RedisConfig
		Email       EmailConfig
		SMTP        SMTPConfig
	}
	PostgresConfig struct {
		Host     string
		Port     string
		Username string
		DBName   string
		Password string
		SSLMode  string
	}
	HTTPConfig struct {
		Host               string        `mapstructure:"host"`
		Port               string        `mapstructure:"port"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMegabytes int           `mapstructure:"maxHeaderBytes"`
	}
	AuthConfig struct {
		JWT          JWTConfig
		PasswordSalt string
	}
	RedisConfig struct {
		Address  string
		Password string
		DB       int
	}
	JWTConfig struct {
		AccessTokenTTL  time.Duration `mapstructure:"accessTokenTTL"`
		RefreshTokenTTL time.Duration `mapstructure:"refreshTokenTTL"`
		SigningKey      string
	}
	EmailConfig struct {
		Templates EmailTemplates
		Subjects  EmailSubjects
	}

	EmailTemplates struct {
		Verification       string `mapstructure:"verification_email"`
		PurchaseSuccessful string `mapstructure:"purchase_successful"`
	}

	EmailSubjects struct {
		Verification       string `mapstructure:"verification_email"`
		PurchaseSuccessful string `mapstructure:"purchase_successful"`
	}

	SMTPConfig struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
		From string `mapstructure:"from"`
		Pass string
	}
)

func Init(configPath string) (*Config, error) {
	populateDefaults()
	if err := parseConfigFile(configPath, os.Getenv("APP_ENV")); err != nil {
		return nil, fmt.Errorf("config.Init: %w", err)
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config.Init: %w", err)
	}

	setFromEnv(&cfg)
	return &cfg, nil

}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("db", &cfg.Postgres); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("auth", &cfg.Auth.JWT); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("redis", &cfg.Redis); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("smtp", &cfg.SMTP); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("email.templates", &cfg.Email.Templates); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("email.subjects", &cfg.Email.Subjects); err != nil {
		return err
	}
	return nil
}

func setFromEnv(cfg *Config) {
	cfg.Postgres.Host = os.Getenv("DB_HOST")
	cfg.Postgres.Port = os.Getenv("DB_PORT")
	cfg.Postgres.Username = os.Getenv("DB_USER")
	cfg.Postgres.DBName = os.Getenv("DB_DBNAME")
	cfg.Postgres.Password = os.Getenv("DB_PASSWORD")
	cfg.Postgres.SSLMode = os.Getenv("DB_SSLMODE")

	cfg.HTTP.Host = os.Getenv("HTTP_HOST")
	cfg.Environment = os.Getenv("APP_ENV")

	cfg.SMTP.Pass = os.Getenv("SMTP_PASSWORD")

	cfg.Redis.Address = os.Getenv("REDIS_URI")
	cfg.Redis.Password = os.Getenv("REDIS_PASSWORD")

	cfg.Auth.PasswordSalt = os.Getenv("PASSWORD_SALT")
	cfg.Auth.JWT.SigningKey = os.Getenv("SIGNING_KEY")
}

func parseConfigFile(folder, env string) error {
	viper.AddConfigPath(folder)
	viper.SetConfigName("main")
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("config.parseConfigFile: %w", err)
	}

	viper.SetConfigName(env)
	return viper.MergeInConfig()

}

func populateDefaults() {
	viper.SetDefault("number.secretCodeTTL", defaultSecretCodeTTL)
	viper.SetDefault("http.port", defaultHTTPPort)
	viper.SetDefault("http.max_header_megabytes", defaultHTTPMaxHeaderMegabytes)
	viper.SetDefault("http.timeouts.read", defaultHTTPRWTimeout)
	viper.SetDefault("http.timeouts.write", defaultHTTPRWTimeout)
	viper.SetDefault("auth.accessTokenTTL", defaultAccessTokenTTL)
	viper.SetDefault("auth.refreshTokenTTL", defaultRefreshTokenTTL)
}
