package agent_manager

import (
	"errors"
	"time"
)

var (
	ErrAgentDoesNotExist   = errors.New("agent does not exist")
	ErrAgentConnEmpty      = errors.New("agent has empty conn string")
	ErrAgentIsNotConnected = errors.New("agent is not connected")
)

type Agent struct {
	ID        string     `db:"id"`
	CreatedAt *time.Time `db:"created_at"`
	Conn      *string    `db:"conn"`
}
