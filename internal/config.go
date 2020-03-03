package internal

import (
	"fmt"
	"os"
	"strings"
)

// options contains all available configs
// for application. All values getting from
// env. For details read .env.example file.
type options struct {
	addr        string
	databaseURL string
}

func newOptions() (options, error) {
	var (
		ok     bool
		opts   options
		missed []string
	)

	if opts.addr, ok = os.LookupEnv("ADDR"); !ok {
		missed = append(missed, "ADDR")
	}
	if opts.databaseURL, ok = os.LookupEnv("DATABASE_URL"); !ok {
		missed = append(missed, "DATABASE_URL")
	}

	if len(missed) != 0 {
		return opts, fmt.Errorf("newOptions fails, variables [%s] doesn't exists", strings.Join(missed, ", "))
	}

	return opts, nil
}
