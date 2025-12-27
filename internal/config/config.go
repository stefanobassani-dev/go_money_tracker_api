package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB     DBConfig
	Server ServerConfig
	Tink   TinkConfig
}

type DBConfig struct {
	host     string
	port     string
	user     string
	password string
	database string
}

type ServerConfig struct {
	Port string
}

type TinkConfig struct {
	BaseUrl      string
	ClientId     string
	ClientSecret string
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

	return &Config{
		DB: DBConfig{
			host:     getEnv("DB_HOST", "localhost"),
			port:     getEnv("DB_PORT", "5432"),
			user:     getEnv("DB_USER", "stef"),
			password: getEnv("DB_PASSWORD", "password"),
			database: getEnv("DB_DATABASE", "money_tracker"),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Tink: TinkConfig{
			BaseUrl:      getEnv("TINK_BASE_URL", ""),
			ClientId:     getEnv("TINK_CLIENT_ID", ""),
			ClientSecret: getEnv("TINK_CLIENT_SECRET", ""),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
