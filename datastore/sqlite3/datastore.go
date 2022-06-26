package sqlite3

import (
	"database/sql"
	"os"

	"github.com/frantjc/sequence/datastore"
)

const Driver = "sqlite3"

func NewDatastore(addr string) (datastore.Datastore, error) {
	f, err := os.Create(addr)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open(Driver, f.Name())
	if err != nil {
		return nil, err
	}

	return &sqlite3Datastore{db}, db.Ping()
}

type sqlite3Datastore struct {
	db *sql.DB
}

func (d *sqlite3Datastore) Ping() error {
	return d.db.Ping()
}

func (d *sqlite3Datastore) Close() error {
	return d.db.Close()
}

func (d *sqlite3Datastore) Driver() string {
	return Driver
}
