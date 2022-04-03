package storage

import (
	"database/sql"

	"github.com/ganzz96/dasha-manager/internal/agents"
	"github.com/pkg/errors"
)

func (db *DB) GetAgent(id string) (*agents.Agent, error) {
	var agent agents.Agent

	err := db.raw.Get(&agent, "SELECT * FROM agents WHERE id = ?", id)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return &agent, errors.WithStack(err)
}

func (db *DB) CreateAgent(agent *agents.Agent) error {
	_, err := db.raw.NamedExec("INSERT INTO agents(id, created_at) VALUES (:id, :created_at)", agent)
	return errors.WithStack(err)
}

func (db *DB) UpdateAgentConn(id string, conn string) error {
	_, err := db.raw.Exec("UPDATE agents SET conn = ? WHERE id = ?", conn, id)
	return errors.WithStack(err)
}
