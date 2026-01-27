package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DB     DBConfig
	Server ServerConfig
	Tink   TinkConfig
	JWT    JWTConfig
	Redis  RedisConfig
}

type DBConfig struct {
	host     string
	port     string
	user     string
	password string
	database string
	Driver   string
}

type ServerConfig struct {
	Port string
}

type TinkConfig struct {
	BaseUrl      string
	ClientId     string
	ClientSecret string
	RedirectUri  string
	Market       string
	Locale       string
}

type JWTConfig struct {
	Secret string
	Expiry time.Duration
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func (c *DBConfig) ConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.user,
		c.password,
		c.host,
		c.port,
		c.database,
	)
}

func Load() *Config {
	_ = godotenv.Load()

	expiryDuration, err := time.ParseDuration(getEnv("JWT_EXPIRY", ""))
	if err != nil {
		log.Printf("Formato JWT_EXPIRY non valido, uso il default di 24h: %v", err)
		expiryDuration = 24 * time.Hour
	}

	return &Config{
		DB: DBConfig{
			host:     getEnv("DB_HOST", "localhost"),
			port:     getEnv("DB_PORT", "5432"),
			user:     getEnv("DB_USER", "stef"),
			password: getEnv("DB_PASSWORD", "password"),
			database: getEnv("DB_DATABASE", "money_tracker"),
			Driver:   getEnv("DB_DRIVER", "postgres"),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Tink: TinkConfig{
			BaseUrl:      getEnv("TINK_BASE_URL", ""),
			ClientId:     getEnv("TINK_CLIENT_ID", ""),
			ClientSecret: getEnv("TINK_CLIENT_SECRET", ""),
			RedirectUri:  getEnv("TINK_REDIRECT_URI", ""),
			Market:       "IT",
			Locale:       "it_IT",
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", ""),
			Expiry: expiryDuration,
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDRESS", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", "password"),
			DB:       0,
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
