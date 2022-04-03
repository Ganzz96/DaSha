package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
)

type Controller struct {
}

func New() *Controller {
	return &Controller{}
}

func (c *Controller) Load(path string) (*AgentConfig, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var cfg AgentConfig

	if err = json.Unmarshal([]byte(file), &cfg); err != nil {
		return nil, errors.WithStack(err)
	}

	return &cfg, nil
}

func (c *Controller) SyncAgentID(agID string) error {
	data, err := json.Marshal(AgentMeta{AgentID: agID})
	if err != nil {
		return errors.WithStack(err)
	}

	path := path.Join(filepath.Dir(os.Args[0]), ".meta")
	if err := ioutil.WriteFile(path, data, os.ModeExclusive); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *Controller) LoadMeta() (*AgentMeta, error) {
	path := path.Join(filepath.Dir(os.Args[0]), ".meta")

	if _, err := os.Stat(".meta"); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, errors.WithStack(err)
	}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var meta AgentMeta

	if err = json.Unmarshal([]byte(file), &meta); err != nil {
		return nil, errors.WithStack(err)
	}

	return &meta, nil
}
