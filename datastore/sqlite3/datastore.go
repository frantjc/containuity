package sqlite3

import (
	"database/sql"

	"github.com/frantjc/sequence/datastore"
)

const Driver = "sqlite3"

func NewDatastore(addr string) (datastore.Datastore, error) {
	db, err := sql.Open(Driver, addr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &sqlite3Datastore{db}, nil
}

type sqlite3Datastore struct {
	db *sql.DB
}

var _ datastore.Datastore = &sqlite3Datastore{}

func (d *sqlite3Datastore) Ping() error {
	return d.db.Ping()
}

func (d *sqlite3Datastore) Close() error {
	return d.db.Close()
}

func (d *sqlite3Datastore) Driver() string {
	return Driver
}
