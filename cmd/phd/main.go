package main

import (
	"ph/internal"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	internal.Run()
}
