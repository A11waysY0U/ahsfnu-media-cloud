package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Upload   UploadConfig
	HMAC     HMACConfig
}

type ServerConfig struct {
	Port string
	Mode string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	SecretKey   string
	ExpireHours int
}

type UploadConfig struct {
	MaxFileSize  int64
	AllowedTypes []string
	UploadPath   string
}

type HMACConfig struct {
	SecretKey string
}

var AppConfig *Config

func Init() {
	// 加载.env文件
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	AppConfig = &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Mode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "ahsfnu_media"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			SecretKey:   getEnv("JWT_SECRET", "your-secret-key"),
			ExpireHours: 24,
		},
		Upload: UploadConfig{
			MaxFileSize:  100 * 1024 * 1024, // 100MB
			AllowedTypes: []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".mp4", ".mov", ".avi", ".wmv", ".flv", ".mkv", ".webm"},
			UploadPath:   getEnv("UPLOAD_PATH", "./uploads"),
		},
		HMAC: HMACConfig{
			SecretKey: getEnv("HMAC_SECRET", "your-hmac-secret-key"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
