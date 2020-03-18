package mi

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v4"
)

const migrationsDir = "./migrations/"
const migrationPostfix = "*.sql"

type Registry struct {
	Migrations []Migration

	conn *pgx.Conn
}

func NewRegistry() *Registry {
	r := Registry{}

	matches, err := filepath.Glob(fmt.Sprint(migrationsDir, migrationPostfix))
	if err != nil {
		log.Fatal(err)
	}

	for _, match := range matches {
		f, err := os.Open(match)
		if err != nil {
			log.Fatal(err)
		}

		reader := bufio.NewReader(f)
		line, _, err := reader.ReadLine()
		if err != nil {
			log.Fatal(err)
		}
		data, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Fatal(err)
		}

		r.Migrations = append(r.Migrations, Migration{
			Name:    match,
			Line:    parseLine(string(line)),
			Content: string(data),
		})
	}

	return &r
}


type Migration struct {
	Name    string
	Line    line
	Content string
}

type line struct {
	Number int
}

const prefix = "-- mi:"

func parseLine(src string) line {
	if !strings.HasPrefix(src, prefix) {
		log.Fatal("line must starts with prefix")
	}

	src = strings.TrimPrefix(src, prefix)
	src = strings.TrimSpace(src)

	number, err := strconv.Atoi(src)
	if err != nil {
		log.Fatal(err)
	}

	return line{
		Number: number,
	}
}
