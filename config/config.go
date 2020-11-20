package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Token string
	Stg   *Settings
}

func NewConfig(settingsFile, token string) (*Config, error) {

	// if a program arg token is empty
	if token == "" {
		// loads values from .env into the system
		if err := godotenv.Load(); err != nil {
			log.Println("no .env file found. Take token from program args")
		}

		//get a token from the env variables
		token = getEnv("TELEGRAM_BOT_TOKEN", "")

		if token == "" {
			return nil, fmt.Errorf("token not found in arg, in .env file, in env variables")
		}
	}

	stg, err := NewSettings(settingsFile)
	if err != nil {
		return nil, err
	}

	return &Config{
		Token: token,
		Stg:   stg,
	}, nil
}

// read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
