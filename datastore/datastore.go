package datastore

type Datastore interface {
	Driver() string
	Ping() error
	Close() error
}
