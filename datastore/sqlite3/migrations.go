package sqlite3

import (
	"database/sql"

	"github.com/frantjc/sequence/datastore"
	"github.com/frantjc/sequence/datastore/sqlite3/migrations"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

type sqlite3Migrations struct {
	migrate *migrate.Migrate
}

func NewMigrations(addr string) (datastore.Migrations, error) {
	db, err := sql.Open(Driver, addr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	instance, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return nil, err
	}

	src, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithInstance(Driver, src, Driver, instance)
	if err != nil {
		return nil, err
	}

	return &sqlite3Migrations{m}, nil
}

func (m *sqlite3Migrations) Up() error {
	return m.migrate.Up()
}

func (m *sqlite3Migrations) Down() error {
	return m.migrate.Down()
}
