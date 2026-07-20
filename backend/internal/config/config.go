package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort            string
	AppEnv             string
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	DBSSLMode          string
	JWTSecret          string
	JWTExpirationHours int
	AdminSecret        string // clave para crear usuarios ADMIN (leída de ADMIN_SECRET)
}

// findEnvFile busca el archivo .env subiendo desde el directorio de trabajo
// hasta encontrarlo o agotar los niveles de búsqueda.
func findEnvFile() string {
	dir, err := os.Getwd()
	if err != nil {
		return ".env"
	}
	for {
		candidate := filepath.Join(dir, ".env")
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			// Llegamos a la raíz del sistema de archivos sin encontrar .env
			break
		}
		dir = parent
	}
	return ".env" // fallback
}

func Load() (*Config, error) {
	envPath := findEnvFile()
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("[WARN] No se pudo cargar el archivo .env desde '%s': %v", envPath, err)
		log.Println("[WARN] Se usarán variables de entorno del sistema o valores por defecto")
	} else {
		log.Printf("[INFO] Variables de entorno cargadas desde: %s", envPath)
	}

	expirationHours, err := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "24"))
	if err != nil {
		return nil, fmt.Errorf("JWT_EXPIRATION_HOURS invalido: %w", err)
	}

	return &Config{
		AppPort:            getEnv("APP_PORT", "8080"),
		AppEnv:             getEnv("APP_ENV", "development"),
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBUser:             getEnv("DB_USER", "postgres"),
		DBPassword:         getEnv("DB_PASSWORD", "postgres"),
		DBName:             getEnv("DB_NAME", "ecommerce"),
		DBSSLMode:          getEnv("DB_SSLMODE", "disable"),
		JWTSecret:          getEnv("JWT_SECRET", "grupo5-super-secret-key-change-in-production"),
		JWTExpirationHours: expirationHours,
		AdminSecret:        getEnv("ADMIN_SECRET", ""), // sin fallback: si no está en .env, nadie puede ser ADMIN
	}, nil
}

func (c *Config) DatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort, c.DBSSLMode,
	)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
