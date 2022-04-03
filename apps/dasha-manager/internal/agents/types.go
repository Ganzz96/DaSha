package agents

import (
	"errors"
	"time"
)

var (
	ErrAgentDoesNotExist = errors.New("agent does not exist")
)

type Agent struct {
	ID        string     `db:"id"`
	CreatedAt *time.Time `db:"created_at"`
	Conn      *string    `db:"conn"`
}

type RegisterRequest struct {
}

type RegisterResponse struct {
	AgentID string `json:"agent_id"`
}
