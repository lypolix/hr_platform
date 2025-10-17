package internal

import (
	"fmt"

	envparser "github.com/caarlos0/env/v11"
)

type environment struct {
	SecretKey   string `env:"SECRET_KEY,required"`
	PostgresDSN string `env:"POSTGRES_DSN,required"`
}

func LoadEnv() (environment, error) {
	var env environment

	err := envparser.Parse(&env)
	if err != nil {
		return environment{}, fmt.Errorf("error parsing config: %w", err)
	}

	return env, nil
}
