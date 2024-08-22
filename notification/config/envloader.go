package config

import (
	"log"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"gopkg.in/go-playground/validator.v9"
)

type OptionsEnv struct {
	Prefix string
	DotEnv bool
}

func envLoader(v any, opt OptionsEnv) (err error) {
	if opt.DotEnv {
		if err = godotenv.Load(); err != nil {
			return err
		}
	}

	err = env.ParseWithOptions(v, env.Options{
		Prefix: opt.Prefix,
	})

	validate := validator.New()
	if err := validate.Struct(v); err != nil {
		log.Fatalf("Validation failed: %v", err)
	}

	return err
}
