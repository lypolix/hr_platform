package internal

import (
	"fmt"

	envparser "github.com/caarlos0/env/v11"
)

type environment struct{}

func LoadEnv() (environment, error) {
	var env environment

	err := envparser.Parse(&env)
	if err != nil {
		return environment{}, fmt.Errorf("error parsing config: %w", err)
	}

	return env, nil
}
