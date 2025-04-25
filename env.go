package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func InitEnv() {
	godotenv.Load()
}

func GetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is not set", key)
		os.Exit(1)
	}
	return value
}
