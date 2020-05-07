// TODO: we need library for config loading
package internal

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// options contains all available configs
// for application. All values getting from
// env. For details read .env.example file.
type options struct {
	addr        string
	databaseURL string

	globalSalt           []byte
	bcryptWorkFactor     int
	minLenForNewPassword int
	tokenTTL             time.Duration
}

func newOptions() (opts options, err error) {
	baseErr := "newOptions fails: %v"

	var ok bool
	var missed []string
	if opts.addr, ok = os.LookupEnv("ADDR"); !ok {
		missed = append(missed, "ADDR")
	}
	if opts.databaseURL, ok = os.LookupEnv("DATABASE_URL"); !ok {
		missed = append(missed, "DATABASE_URL")
	}
	globalSalt, ok := os.LookupEnv("GLOBAL_SALT")
	if !ok {
		missed = append(missed, "GLOBAL_SALT")
	}
	if len(globalSalt) < 32 {
		return opts, fmt.Errorf(baseErr, "GLOBAL_SALT length MUST be grates or equal 32")
	}
	opts.globalSalt = []byte(globalSalt)
	opts.bcryptWorkFactor = 10
	passwordLen, ok := os.LookupEnv("MIN_LEN_FOR_NEW_PASSWORD")
	if !ok {
		missed = append(missed, "MIN_LEN_FOR_NEW_PASSWORD")
	}
	opts.minLenForNewPassword, err = strconv.Atoi(passwordLen)
	if err != nil {
		return opts, fmt.Errorf(baseErr, fmt.Sprintf("minLenForNewPassword read fails: %v", err))
	}
	opts.tokenTTL, err = time.ParseDuration(os.Getenv("TOKEN_TTL"))
	if err != nil {
		return opts, fmt.Errorf(baseErr, fmt.Sprintf("TOKEN_TTL read fails: %v", err))
	}

	if len(missed) != 0 {
		return opts, fmt.Errorf("newOptions fails, variables [%s] doesn't exists", strings.Join(missed, ", "))
	}

	return opts, nil
}
