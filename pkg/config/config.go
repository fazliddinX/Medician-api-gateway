package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"

	_ "github.com/lib/pq"
)

type Config struct {
	SECRET_KEY_ACCESS  string
	GIN_SERVER_PORT    string
	USER_SERVER_PORT   string
	HEALTH_SERVER_PORT string
	RABBIT_URL         string

	RADIS string

	USER_HOST string

	DB_PORT     string
	DB_USER     string
	DB_NAME     string
	DB_PASSWORD string
	DB_HOST     string

	HEALTH_HOST string
	REDIS_HOST  string
}

func Load() Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found")
	}

	config := Config{}

	config.GIN_SERVER_PORT = cast.ToString(coalesce("GIN_SERVER_PORT", ":8081"))
	config.USER_SERVER_PORT = cast.ToString(coalesce("USER_SERVER_PORT", ":50050"))
	config.HEALTH_SERVER_PORT = cast.ToString(coalesce("HEALTH_SERVER_PORT", ":50051"))

	config.USER_HOST = cast.ToString(coalesce("USER_HOST", "medical-auth"))
	config.HEALTH_HOST = cast.ToString(coalesce("HEALTH_HOST", "medical-auth"))

	config.SECRET_KEY_ACCESS = cast.ToString(coalesce("SECRET_KEY_ACCESS", "secret_key"))
	config.RABBIT_URL = cast.ToString(coalesce("RABBIT_URL", "amqp://guest:guest@rabbit:5672/"))
	config.DB_PORT = cast.ToString(coalesce("DB_PORT", "5432"))
	config.DB_USER = cast.ToString(coalesce("DB_USER", "casbin"))
	config.DB_NAME = cast.ToString(coalesce("DB_NAME", "casbin"))
	config.DB_PASSWORD = cast.ToString(coalesce("DB_PASSWORD", "123321"))
	config.REDIS_HOST = cast.ToString(coalesce("REDIS_HOST", "redis"))
	config.DB_HOST = cast.ToString(coalesce("DB_HOST", "postgres"))
	config.RADIS = cast.ToString(coalesce("RADIS", "RADIS"))

	return config
}

func coalesce(env string, defaultValue interface{}) interface{} {
	value, exists := os.LookupEnv(env)
	if !exists {
		return defaultValue
	}
	return value
}
