package internal

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

// Opts contains all available configs
// for application. All values getting from
// env. For details read .env.example file.
type Opts struct {
	Addr                 string        `yaml:"Addr"`
	DatabaseURL          string        `yaml:"DatabaseURL"`
	GlobalSalt           string        `yaml:"GlobalSalt"`
	BcryptWorkFactor     int           `yaml:"BcryptWorkFactor"`
	MinLenForNewPassword int           `yaml:"MinLenForNewPassword"`
	TokenTTL             time.Duration `yaml:"TokenTTL"`
}

func FromFile(configFilename string) (opts Opts, err error) {
	buf, err := ioutil.ReadFile(configFilename)
	if err != nil {
		return Opts{}, err
	}
	if err := yaml.UnmarshalStrict(buf, &opts); err != nil {
		return Opts{}, err
	}

	var missed []string
	if opts.Addr == "" {
		missed = append(missed, "Addr")
	}
	if opts.DatabaseURL == "" {
		missed = append(missed, "DatabaseURL")
	}
	if len(opts.GlobalSalt) < 32 {
		return opts, errors.New("GlobalSalt length MUST be grates or equal 32")
	}
	if opts.BcryptWorkFactor == 0 {
		missed = append(missed, "BcryptWorkFactor")
	}
	if opts.MinLenForNewPassword == 0 {
		missed = append(missed, "MinLenForNewPassword")
	}
	if opts.TokenTTL == 0 {
		missed = append(missed, "TokenTTL")
	}

	if len(missed) != 0 {
		return opts, fmt.Errorf("FromFile fails, [%s] fields doesn't exist", strings.Join(missed, ", "))
	}

	return opts, err
}
