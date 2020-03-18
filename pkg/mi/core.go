package mi

import (
	"errors"
	"fmt"
)

const migrationNameTemplate = "%03d.%s.sql"

func formatMigrationFileName(number int, name string) (string, error) {
	if number < 0 {
		return "", errors.New("migration number must be great or equal 0")
	}
	if name == "" {
		return "", errors.New("name not passed")
	}
	if number > 999 {
		return "", errors.New("number too big")
	}
	return fmt.Sprintf(migrationNameTemplate, number, name), nil
}

func parseMigrationFileName() {

}
