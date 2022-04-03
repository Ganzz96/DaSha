package storage

import (
	"github.com/ganzz96/dasha-manager/internal/agents"
	"github.com/pkg/errors"
)

func (db *DB) UpdateAgentConn(id string, conn string) error {
	_, err := db.raw.Exec("UPDATE agents SET conn = ? WHERE id = ?", conn, id)
	return errors.WithStack(err)
}

func (db *DB) GetAgent(id string) (*agents.Agent, error) {
	var agent agents.Agent

	err := db.raw.Get(&agent, "SELECT * FROM agents WHERE id = ?", id)
	return &agent, errors.WithStack(err)
}
