package main

import (
	"os"

	"github.com/joho/godotenv"
	"ph/internal"
)

func main() {
	file := os.Getenv("ENV")
	if file == "" {
		file = ".env"
	}
	_ = godotenv.Load(file)
	internal.Run()
}
