package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	AppEnv     string `env:"APP_ENV" default:"dev"`
	Debug      bool   `env:"DEBUG" default:"true"`
	DBName     string `env:"DB_NAME" required:"true"`
	DBUser     string `env:"DB_USER" required:"true"`
	DBPassword string `env:"DB_PASSWORD" required:"true"`
	DBHost     string `env:"DB_HOST" required:"true"`
	DBPort     string `env:"DB_PORT" required:"true"`

	RedisPort     string `env:"REDIS_PORT" required:"true"`
	RedisHost     string `env:"REDIS_HOST" required:"true"`
	RedisPassword string `env:"REDIS_PASSWORD" required:"false"`
	RedisDB       int    `env:"REDIS_DB" required:"true"`

	GrpcServerPort string `env:"GRPC_SERVER_PORT" required:"true"`

	ElasticHost string `env:"ELASTIC_HOST" required:"true"`
	ElasticPort string `env:"ELASTIC_PORT" required:"true"`
}

func Load() *Config {
	var cfg Config
	err := cleanenv.ReadConfig("./.env", &cfg)
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}
	return &cfg
}
