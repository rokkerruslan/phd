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
	opts, err := FromEnv()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("OPTS: %+v\n", opts)

	ctx, cancel := context.WithTimeout(context.Background(), 42*time.Second)
	defer cancel()

	m := NewMigrator(ctx, opts)

	if err := m.Init(ctx); err != nil {
		fmt.Println(err)
		return
	}

	args := os.Args[1:]
	if len(args) == 0 {
		m.Status()
		return
	}
	command := args[0]

	switch command {
	case "n", "ne", "new":
		if len(args) < 2 {
			log.Fatal("you need set name for migration")
		}
		m.New(args[1], opts)
	case "u", "up":
		// todo rr: to number?
		n := -1
		m.Migrate(n, opts)
	case "s", "st", "sta", "statu", "status":
		m.Status()
	case "h", "he", "hel", "help":
		fmt.Print(Usage)
	default:
		fmt.Println("undefined command:", command)
		fmt.Println(Usage)
	}
}
