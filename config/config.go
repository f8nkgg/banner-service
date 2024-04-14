package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type (
	Config struct {
		App        `yaml:"app"`
		HTTPServer `yaml:"HTTPServer"`
		Log        `yaml:"logger"`
		PG         `yaml:"postgres"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTPServer struct {
		Host            string        `yaml:"host" env:"SERVER_HOST" env-default:"localhost"`
		Port            string        `env-required:"true" yaml:"port" env:"SERVER_PORT"`
		ReadTimeout     time.Duration `yaml:"read_timeout" env:"READ_TIMEOUT" env-default:"5s"`
		WriteTimeout    time.Duration `yaml:"write_timeout" env:"WRITE_TIMEOUT" env-default:"5s"`
		MaxHeaderBytes  int           `yaml:"max_header_bytes" env:"MAX_HEADER_BYTES" env-default:"1048576"`
		ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env:"SHUTDOWN_TIMEOUT" env-default:"3s"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}

	PG struct {
		PoolMax      int           `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
		ConnAttempts int           `yaml:"conn_attempts" env-default:"5"`
		ConnTimeout  time.Duration `yaml:"conn_timeout" env-default:"10s"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
