package config

import (
	"github.com/caarlos0/env"
)

const ServiceName = "pikachu"

type cfg struct {
	Env          string `env:"ENV" envDefault:"dev"`
	BasePath     string `env:"BASE_PATH"`
	Host         string `env:"HOST" envDefault:"http://localhost:3000"`
	Port         int    `env:"PORT" envDefault:"80"`
	InternalPort int    `env:"INTERNAL_PORT" envDefault:"3000"`
	LogLevel     string `env:"LOG_LEVEL" envDefault:"debug"`
	DBHostname   string `env:"DB_HOSTNAME" envDefault:"127.0.0.1"`
	DBPort       int    `env:"DB_PORT" envDefault:"3306"`
	DBUsername   string `env:"DB_USERNAME" envDefault:"root"`
	DBPassword   string `env:"DB_PASSWORD" envDefault:"123456"`
	DBDatabase   string `env:"DB_DATABASE" envDefault:"pikachu"`
	DDTrace      bool   `env:"DD_TRACE" envDefault:"true"`

	SentryDSN      string   `env:"SENTRY_DSN"`
	KafkaHost      []string `env:"KAFKA_HOST"`
	RedisCacheHost string   `env:"REDIS_CACHE_HOST" envDefault:"127.0.0.1"`
	RedisCachePort int      `env:"REDIS_CACHE_PORT" envDefault:"6379"`
}

var Cfg *cfg

func InitConfig() {
	Cfg = &cfg{}

	if err := env.Parse(Cfg); err != nil {
		panic(err)
	}

}
