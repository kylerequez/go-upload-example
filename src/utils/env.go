package utils

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	return nil
}

func GetEnv(key string) string {
	v, isExists := os.LookupEnv(key)
	if !isExists {
		return ""
	}

	return v
}
