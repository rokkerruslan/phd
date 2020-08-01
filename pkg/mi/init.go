package main

import (
	"context"
)

const migrationsTableTemplate = `
CREATE TABLE IF NOT EXISTS public.migrations (
    id integer NOT NULL,
    created timestamp without time zone,
    name text,
    number integer
);
`

func (m *Migrator) Init(ctx context.Context) (error) {
	// todo rr: check structure?
	_, err := m.db.Exec(ctx, migrationsTableTemplate)
	return err
}
