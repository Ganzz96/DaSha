package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
)

type UagentConfig struct {
	HostPort             string `json:"host_port"`
	StunServer           string `json:"stun_server"`
	AgentManagerEndpoint string `json:"agent_manager_endpoint"`
}

func Load(path string) (*UagentConfig, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var cfg UagentConfig

	if err = json.Unmarshal([]byte(file), &cfg); err != nil {
		return nil, errors.WithStack(err)
	}

	return &cfg, nil
}
