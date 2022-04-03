package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
)

type ClientConfig struct {
	HostPort string `json:"host_port"`
}

func Load(path string) (*ClientConfig, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var cfg ClientConfig

	if err = json.Unmarshal([]byte(file), &cfg); err != nil {
		return nil, errors.WithStack(err)
	}

	return &cfg, nil
}
