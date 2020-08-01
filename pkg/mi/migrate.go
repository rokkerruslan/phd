package main

import (
	"context"
	"fmt"
	"log"
	"sort"
	"time"
)

// todo rr: --fake option?
func (m *Migrator) Migrate(to int, opts Opts) {
	r := NewRegistry(opts, m.db)
	r.Sort()
	r.Apply()
	r.Commit()
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
