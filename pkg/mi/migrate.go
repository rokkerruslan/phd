package mi

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/jackc/pgx/v4"
)

func Migrate() {
	r := NewRegistry()
	r.InitConnection()
	r.CheckMigrationTable("public", "migrations")
	r.Sort()
	r.Apply()
	r.Commit()
}

func (r *Registry) InitConnection() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	r.conn = conn
}

var existQuery = `
SELECT EXISTS (
	SELECT FROM information_schema.tables
	WHERE  table_schema = $1
	AND    table_name   = $2
);
`

var createMigrationTableQuery = `
CREATE TABLE migrations (
	id serial PRIMARY KEY,
	created timestamp without time zone,
	name text,
	number integer
);
`

func (r *Registry) CheckMigrationTable(schema, table string) {
	var isExist bool
	if err := r.conn.QueryRow(context.Background(), existQuery, schema, table).Scan(&isExist); err != nil {
		log.Fatal(err)
	}
	if !isExist {
		if _, err := r.conn.Exec(context.Background(), createMigrationTableQuery); err != nil {
			log.Fatal(err)
		}
	}
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
	rows, err := r.conn.Query(context.Background(), "SELECT id, created, name, number FROM public.migrations")
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

func (r *Registry) Commit() {
	for _, m := range r.Migrations {
		_, err := r.conn.Exec(
			context.Background(),
			"INSERT INTO public.migrations (created, name, number) VALUES (NOW(), $1, $2)",
			m.Name,
			m.Line.Number,
		)
		if err != nil {
			log.Fatal(err)
		}
	}
}
