package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App  `yaml:"app"`
		HTTP `yaml:"http"`
		Log  `yaml:"logger"`
		Dapr `yaml:"dapr"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Port             string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
		NodeJsRuntimeUrl string `env-required:"true" yaml:"nodeJsRuntimeUrl" env:"HTTP_NODEJS_RUNTIME_URL"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"logLevel"   env:"LOG_LEVEL"`
	}

	Dapr struct {
		StoreName string `env-required:"true" yaml:"storeName"   env:"DAPR_STORE_NAME"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./internal/config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
