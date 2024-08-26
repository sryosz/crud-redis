package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Http httpConfig
	DB   dbConfig
}

type httpConfig struct {
	Port uint16 `yaml:"http_port" env-default:"8000"`
}

type dbConfig struct {
	Host            string `env:"DB_HOST"               env-default:"localhost"`
	Port            string `env:"DB_PORT"               env-default:"5432"`
	User            string `env:"DB_USER"               env-default:"postgres"`
	Pass            string `env:"DB_PASS"               env-default:"root"`
	Name            string `env:"DB_NAME"               env-default:"redis"`
	MaxConns        int32  `env:"DB_MAX_CONNS"          env-default:"10"`
	MinConns        int32  `env:"DB_MIN_CONNS"          env-default:"1"`
	MaxConnLifetime int64  `env:"DB_MAX_CONN_LIFETIME"  env-default:"3600000000000"`
	MaxConnIdleTime int64  `env:"DB_MAX_CONN_IDLE_TIME" env-default:"1800000000000"`
}

func MustLoad() *Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("Failed to read config: " + err.Error())
	}

	return &cfg
}
