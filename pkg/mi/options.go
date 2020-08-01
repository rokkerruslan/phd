package main

import (
	"errors"
	"os"
	"path"
)

const DefaultMigrationsDir = "./migrations/"

type Opts struct {
	DatabaseURL       string
	MigrationsPattern string
}

func FromEnv() (opts Opts, err error) {
	var ok bool
	opts.DatabaseURL, ok = os.LookupEnv("DATABASE_URL")
	if !ok {
		return opts, errors.New("DATABASE_URL not found")
	}

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	opts.MigrationsPattern = path.Join(dir, DefaultMigrationsDir, "*.sql")

	return opts, nil
}
