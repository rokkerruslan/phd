package mi

import (
	"fmt"
	"log"
	"os"
)

const Usage = `
$ mi -h
`

func App() {
	args := os.Args[1:]

	if len(args) == 0 {
		Status()
		return
	}

	log.SetFlags(log.LstdFlags | log.Llongfile)

	command := args[0]

	switch command {
	case "n", "ne", "new":
		if len(args) < 2 {
			log.Fatal("you need set name for migration")
		}
		New(args[1])
	case "u", "up":
		Migrate()
	case "s", "st", "sta", "statu", "status":
		Status()
	case "h", "he", "hel", "help":
		fmt.Print(Usage)
	}
}
