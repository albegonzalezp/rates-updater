package config

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

const (
	EnvDev     = "dev"
	EnvStaging = "staging"
	EnvProd    = "prod"
)

func LoadEnvironmentVariables() error {
	// Detect environment
	env := os.Getenv("ENV")

	// can be dev, staging, prod
	switch env {
	case EnvProd:
		log.Println("Loading envs for ENV: prod")

		if err := checkEnvVariables(env); err != nil {
			return err
		}

	case EnvStaging:

		log.Println("Loading envs for ENV: staging")

		if err := checkEnvVariables(env); err != nil {
			return err
		}

	default:

		env = EnvDev

		if err := os.Setenv("ENV", env); err != nil {
			return err
		}

		// Fallback to dev
		log.Println("Loading envs for ENV: dev")

		// Load from local.
		if err := godotenv.Load(".env"); err != nil {
			return err
		}

		if err := checkEnvVariables(env); err != nil {
			return err
		}
	}

	return nil
}

func checkEnvVariables(env string) error {

	if os.Getenv("EXCHANGE_API_KEY") == "" ||
		os.Getenv("EXCHANGE_API_URL") == "" ||
		os.Getenv(fmt.Sprintf("DB_HOST_%s", strings.ToUpper(env))) == "" ||
		os.Getenv(fmt.Sprintf("DB_USER_%s", strings.ToUpper(env))) == "" ||
		os.Getenv(fmt.Sprintf("DB_PASSWORD_%s", strings.ToUpper(env))) == "" ||
		os.Getenv(fmt.Sprintf("DB_NAME_%s", strings.ToUpper(env))) == "" ||
		os.Getenv(fmt.Sprintf("DB_PORT_%s", strings.ToUpper(env))) == "" {
		return errors.New(fmt.Sprintf("env variables for %s are not set", env))
	}
	return nil
}
