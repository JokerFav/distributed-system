package config

import "os"

type Config struct {
	HTTPAddress string
	DatabaseURL string
}

func Load() Config {
	return Config{
		HTTPAddress: os.Getenv("HTTP_ADDRESS"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
}
