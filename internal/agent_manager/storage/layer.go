package storage

import (
	"database/sql"

	"github.com/ganzz96/dasha/internal/agent_manager"
	"github.com/pkg/errors"
)

func (db *DB) GetAgent(id string) (*agent_manager.Agent, error) {
	var agent agent_manager.Agent

	err := db.raw.Get(&agent, "SELECT * FROM agents WHERE id = ?", id)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return &agent, errors.WithStack(err)
}

func (db *DB) CreateAgent(agent *agent_manager.Agent) error {
	_, err := db.raw.NamedExec("INSERT INTO agents(id, created_at) VALUES (:id, :created_at)", agent)
	return errors.WithStack(err)
}

func (db *DB) UpdateAgentConn(id string, conn string) error {
	_, err := db.raw.Exec("UPDATE agents SET conn = ? WHERE id = ?", conn, id)
	return errors.WithStack(err)
}
