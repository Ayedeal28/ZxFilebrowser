package config

import "os"

type Config struct {
	RootDir string
	Port    string
}

var AppConfig Config

func Init() {
	AppConfig = Config{
		RootDir: "C:\\Users\\ayede\\Desktop\\TestFile",
		Port:    getEnv("PORT", ":8080"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
