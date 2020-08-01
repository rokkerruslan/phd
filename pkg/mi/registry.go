package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jackc/pgx/v4"
)

type Registry struct {
	Migrations []Migration

	conn *pgx.Conn
}

// todo rr: check inconsistency?
func NewRegistry(opts Opts, conn *pgx.Conn) *Registry {
	r := Registry{
		conn: conn,
	}

	matches, err := filepath.Glob(opts.MigrationsPattern)
	if err != nil {
		log.Fatal(err)
	}

	for _, match := range matches {
		f, err := os.Open(match)
		if err != nil {
			log.Fatal(err)
		}

		number, name, err := parseMigrationFileName(filepath.Base(match))
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
			Name:    name,
			Number:  number,
			Line:    parseLine(string(line)),
			Content: string(data),
		})
	}

	return &r
}

type Migration struct {
	Name    string
	Number  int
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

	number, _, err := parseMigrationFileName(src)
	if err != nil {
		log.Fatal(err)
	}

	return line{
		Number: number,
	}
}
