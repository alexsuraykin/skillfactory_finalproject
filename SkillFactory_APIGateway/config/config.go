package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	APIPort  string `default:"8881"`
	LogLevel string `default:"debug"`
}

func InitConfig() (*Config, error) {
	var cnf Config

	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}

	cnf = Config{
		APIPort: os.Getenv("SERVER_PORT"),
	}

	return &cnf, nil
}
