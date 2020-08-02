package main

import (
	"fmt"
	"os"

	"ph/internal"
)

func main() {
	args := os.Args[1:]

	if len(args) != 1 {
		fmt.Println("usage: phd config-file.yaml")
		os.Exit(1)
	}

	internal.Run(args[0])
}
