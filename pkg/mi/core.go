package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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

func parseMigrationFileName(name string) (int, string, error) {
	baseErr := "parseMigrationFileName fails: %v"

	numberIndex := strings.Index(name, ".")
	if numberIndex == -1 {
		return 0, "", errors.New("name don't contains delimiter")
	}
	numberSub := name[:numberIndex]
	if numberSub == "" {
		return 0, "", errors.New("number is empty")
	}
	number, err := strconv.Atoi(numberSub)
	if err != nil {
		return 0, "", fmt.Errorf(baseErr, err)
	}

	rest := name[numberIndex+1:]
	if rest == "" {
		return 0, "", fmt.Errorf(baseErr, "name is empty")
	}
	nameIndex := strings.Index(rest, ".")
	if numberIndex == -1 {
		return 0, "", fmt.Errorf(baseErr, "name don't contains delimiter")
	}
	nameSub := rest[:nameIndex]
	if nameSub == "" {
		return 0, "", fmt.Errorf(baseErr, "name is empty")
	}

	return number, nameSub, nil
}
