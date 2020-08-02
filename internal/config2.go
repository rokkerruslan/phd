package internal

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

// todo rr: drop config and rename this file

// Opts contains all available configs
// for application. All values getting from
// env. For details read .env.example file.
type Opts struct {
	Addr        string `yaml:"Addr"`
	DatabaseURL string `yaml:"DatabaseURL"`

	GlobalSalt           string        `yaml:"GlobalSalt"`
	BcryptWorkFactor     int           `yaml:"BcryptWorkFactor"`
	MinLenForNewPassword int           `yaml:"MinLenForNewPassword"`
	TokenTTL             time.Duration `yaml:"TokenTTL"`
}

func FromFile() (opts Opts, err error) {
	buf, err := ioutil.ReadFile("example.opts.yml")
	if err != nil {
		return Opts{}, err
	}
	if err := yaml.UnmarshalStrict(buf, &opts); err != nil {
		return Opts{}, err
	}

	opts.BcryptWorkFactor = 10

	return opts, err
}
