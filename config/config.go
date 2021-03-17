package config

import (
	"fmt"

	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Config describes the structure of imported config
type Config struct {
	API struct {
		Apikey   string `envconfig:"API_KEY" validate:"required"`
		Endpoint string `envconfig:"API_ENDPOINT" validate:"required"`
	}
	Cache struct {
		Duration int `envconfig:"CACHE_DURATION" validate:"required,gt=0"`
	}
	Database struct {
		Driver string `envconfig:"DB_DRIVER" validate:"required"`
		Host   string `envconfig:"DB_HOST" validate:"required"`
		Port   int    `envconfig:"DB_PORT" validate:"required"`
	}
	Server struct {
		Host string `envconfig:"SERVER_HOST" validate:"required"`
		Port int    `envconfig:"SERVER_PORT" validate:"required"`
	}
}

// LoadConfig reads the configuration from env
func LoadConfig(validate *validator.Validate, cfg *Config) (err error) {
	err = godotenv.Load()
	if err != nil {
		return fmt.Errorf("Error loading .env file:\n %v", err)
	}

	err = envconfig.Process("", cfg)
	if err != nil {
		return fmt.Errorf("Error parsing .env file:\n %v", err)
	}

	err = validate.Struct(cfg)
	if err != nil {
		return fmt.Errorf("Error validating config:\n %v", err)
	}

	return nil
}
