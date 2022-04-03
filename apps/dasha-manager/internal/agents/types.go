package agents

import "time"

type Agent struct {
	ID        string     `db:"id"`
	CreatedAt *time.Time `db:"created_at"`
	Conn      string     `db:"conn"`
}
