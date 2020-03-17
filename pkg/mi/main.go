// TODO: initialize with migrations table
// TODO: console ui
package mi

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
)

const prefix = "-- mi:"

type Migration struct {
	Name    string
	Line    line
	Content string
}

type line struct {
	Number int
}

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

type Registry struct {
	Migrations []Migration

	conn *pgx.Conn
}

func (r *Registry) Sort() {
	sort.Slice(r.Migrations, func(i, j int) bool {
		return r.Migrations[i].Line.Number < r.Migrations[j].Line.Number
	})
}

func (r *Registry) Apply() {
	tx, err := r.conn.Begin(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Filter applied
	r.Filter()

	// Apply all migrations
	for _, migration := range r.Migrations {
		fmt.Println("APPLY: ", migration)
		if _, err := r.conn.Exec(context.Background(), migration.Content); err != nil {
			log.Fatal(err)
		}
	}

	// Commit tx
	if err := tx.Commit(context.Background()); err != nil {
		log.Fatal(err)
	}
}

type MLine struct {
	ID      int
	Created time.Time
	Name    string
	Number  int
}

func (r *Registry) Filter() {
	rows, err := r.conn.Query(context.Background(), "SELECT id, created, name, number FROM infra.migrations")
	if err != nil {
		log.Fatal(err)
	}

	var applied = make(map[int]bool)
	for rows.Next() {
		var line MLine
		if err := rows.Scan(&line.ID, &line.Created, &line.Name, &line.Number); err != nil {
			log.Fatal(err)
		}
		fmt.Println("APPLIED:", line)
		applied[line.Number] = true
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	var nonApplied []Migration
	for _, m := range r.Migrations {
		ok := applied[m.Line.Number]
		if !ok {
			nonApplied = append(nonApplied, m)
		}
	}

	r.Migrations = nonApplied
}

func Migrate() {
	matches, err := filepath.Glob("./migrations/*.sql")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Unsorted Migrations:", matches)

	registry := Registry{}

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

		registry.Migrations = append(registry.Migrations, Migration{
			Name:    match,
			Line:    parseLine(string(line)),
			Content: string(data),
		})
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	registry.conn = conn
	defer conn.Close(context.Background())

	registry.Sort()
	registry.Apply()
	registry.Commit()
}

func (r *Registry) Commit() {
	for _, m := range r.Migrations {
		_, err := r.conn.Exec(
			context.Background(),
			"INSERT INTO infra.migrations (created, name, number) VALUES (NOW(), $1, $2)",
			m.Name,
			m.Line.Number,
		)
		if err != nil {
			log.Fatal(err)
		}
	}
}
