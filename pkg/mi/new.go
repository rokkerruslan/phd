package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func (m *Migrator) New(name string, opts Opts)  {
	r := NewRegistry(opts, m.db)

	migrations := r.Migrations[len(r.Migrations)-1]

	nextMigrationNumber := migrations.Line.Number + 1

	fileName, err := formatMigrationFileName(nextMigrationNumber, name)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(filepath.Join(DefaultMigrationsDir, fileName))
	if err != nil {
		log.Fatal(err)
	}

	if _, err := f.WriteString(fmt.Sprintf(migrationTemplate, fileName)); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Success, your migration:", fileName)
}
