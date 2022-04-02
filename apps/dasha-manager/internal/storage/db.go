package storage

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"

	"github.com/ganzz96/dasha-manager/internal/log"
)

type DB struct {
	logger *log.Logger
	raw    *sqlx.DB
}

func New(logger *log.Logger, path string) (*DB, error) {
	raw, err := sqlx.Connect("sqlite3", path)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	db := &DB{raw: raw, logger: logger}
	return db, db.applyMigrations(path)
}

func (db *DB) applyMigrations(dbPath string) error {
	migrator, err := migrate.New("file://internal/storage/migrations", fmt.Sprintf("sqlite://%s", dbPath))
	if err != nil {
		return errors.WithStack(err)
	}

	migrator.Log = db.logger

	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return errors.WithStack(err)
	}

	return nil
}
