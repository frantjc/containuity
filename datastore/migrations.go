package datastore

type Migrations interface {
	Up() error
	Down() error
}
