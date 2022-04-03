package storage

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	sqllog "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zapadapter"

	"github.com/ganzz96/dasha-manager/internal/log"
)

type DB struct {
	logger *log.Logger
	raw    *sqlx.DB
}

func New(logger *log.Logger, dsn string) (*DB, error) {
	sqldb, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	sqldb = sqllog.OpenDriver(dsn, sqldb.Driver(), zapadapter.New(logger.Logger))
	if err := sqldb.Ping(); err != nil {
		sqldb.Close()
		return nil, err
	}

	db := &DB{raw: sqlx.NewDb(sqldb, "sqlite3"), logger: logger}
	return db, db.applyMigrations(dsn)
}

func (db *DB) applyMigrations(dsn string) error {
	migrator, err := migrate.New("file://internal/storage/migrations", fmt.Sprintf("sqlite://%s", dsn))
	if err != nil {
		return errors.WithStack(err)
	}

	migrator.Log = db.logger

	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return errors.WithStack(err)
	}

	return nil
}
