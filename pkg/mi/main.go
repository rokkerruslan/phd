package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

const Usage = `
Usage: mi -h
`

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		Status()
		return
	}

	command := args[0]

	opts, err := FromEnv()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("OPTS: %+v\n", opts)

	ctx, cancel := context.WithTimeout(context.Background(), 42*time.Second)
	defer cancel()

	m := NewMigrator(ctx, opts)

	switch command {
	case "i", "in", "ini", "init":
		if len(args) != 1 {
			fmt.Println("init no takes parameters")
			return
		}
		if err := m.Init(ctx); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("done")
	case "n", "ne", "new":
		if len(args) < 2 {
			log.Fatal("you need set name for migration")
		}
		New(args[1], opts)
	case "u", "up":
		Migrate(opts)
	case "s", "st", "sta", "statu", "status":
		Status()
	case "h", "he", "hel", "help":
		fmt.Print(Usage)
	default:
		fmt.Println("undefined command:", command)
		fmt.Println(Usage)
	}
}
