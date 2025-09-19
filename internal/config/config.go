package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	AppName string
	AppEnv  string
	AppPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	JWTSecret  string

	Logger *logrus.Logger
}

func LoadConfig() *Config {
	// Load .env file if present (safe to ignore in prod)
	_ = godotenv.Load()

	// Configure logger
	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	appEnv := getEnv("APP_ENV", "development")

	if appEnv == "production" {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
		logger.SetLevel(logrus.InfoLevel)
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
			PadLevelText:    true,
		})
		logger.SetLevel(logrus.DebugLevel)
	}

	cfg := &Config{
		AppName:    getEnv("APP_NAME", "myapp"),
		AppEnv:     appEnv,
		AppPort:    getEnv("APP_PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "myapp_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		JWTSecret:  getEnv("JWT_SECRET", "super-secret-key"),
		Logger:     logger,
	}

	cfg.Logger.Debugf("Config loaded: env=%s port=%s db=%s", cfg.AppEnv, cfg.AppPort, cfg.DBName)
	return cfg
}

// DatabaseDSN returns a Postgres DSN string (for GORM or pgx)
func (c *Config) DatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort, c.DBSSLMode,
	)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
