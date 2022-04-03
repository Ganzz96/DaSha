package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
)

type ManagerConfig struct {
	DBPath       string `json:"db_path"`
	HostPort     string `json:"host_port"`
	AgentMonitor struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"agent_monitor"`
}

func Load(path string) (*ManagerConfig, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var cfg ManagerConfig

	if err = json.Unmarshal([]byte(file), &cfg); err != nil {
		return nil, errors.WithStack(err)
	}

	return &cfg, nil
}
