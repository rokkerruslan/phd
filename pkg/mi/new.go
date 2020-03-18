package mi

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const migrationTemplate = `-- mi: %d

-- WRITE YOUR MIGRATION HERE

`

func New(name string) {
	r := NewRegistry()

	m := r.Migrations[len(r.Migrations) - 1]

	nextMigrationNumber := m.Line.Number + 1

	fileName := fmt.Sprintf("%d.sql", nextMigrationNumber)

	f, err := os.Create(filepath.Join(migrationsDir, fileName))
	if err != nil {
		log.Fatal(err)
	}

	if _, err := f.WriteString(fmt.Sprintf(migrationTemplate, nextMigrationNumber)); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("|%-6s|%-6s|\n", "foo", "b")

	fmt.Println("Success, your migration:", fileName)
}
