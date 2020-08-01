package main

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFormatMigrationFileName_Positive(t *testing.T) {
	got, err := formatMigrationFileName(1, "Positive")
	if err != nil {
		t.Fatalf("expected name, got error: %v", err)
	}
	expected := "001.Positive.sql"

	if diff := cmp.Diff(got, expected); diff != "" {
		t.Error(diff)
	}
}

func TestFormatMigrationFileName_EmptyName(t *testing.T) {
	got, err := formatMigrationFileName(0, "")
	if err == nil {
		t.Fatalf("expected error, got %s", got)
	}
	if !strings.Contains(err.Error(), "name") {
		t.Errorf("expected error with name, got %v", err)
	}
}

func TestFormatMigrationFileName_BigNumber(t *testing.T) {
	got, err := formatMigrationFileName(1000, "BigNumber")
	if err == nil {
		t.Fatalf("expected error, got %s", got)
	}
	if !strings.Contains(err.Error(), "number") {
		t.Errorf("expected error with number, got %v", err)
	}
}

func TestParseMigrationFileName_Positive(t *testing.T) {
	number, name, err := parseMigrationFileName("001.Positive.sql")
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(1, number); diff != "" {
		t.Error(diff)
	}
	if diff := cmp.Diff("Positive", name); diff != "" {
		t.Error(diff)
	}
}
