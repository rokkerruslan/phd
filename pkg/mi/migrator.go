package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
)

type Migrator struct {
	opts Opts

	db *pgx.Conn
}

func NewMigrator(ctx context.Context, opts Opts) *Migrator {
	var m Migrator

	var err error
	m.db, err = pgx.Connect(ctx, opts.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	return &m
}
