package config

import (
	"os"
	"strconv"
)

type Config struct {
	TCPAddr  string
	GameSize int
	Mines    int
}

func Load() Config {
	return Config{
		TCPAddr:  getEnv("TCP_ADDR", ":8080"),
		GameSize: getEnvInt("GAME_SIZE", 8),
		Mines:    getEnvInt("MINES", 10),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	n, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return n
}
