package store

import (
	"tirelease/utils/configs"
	"tirelease/utils/database"
)

var tempDB *database.Database

func NewDBClients(dsn string) *database.Database {
	db := database.MustConnect(dsn)
	return db
}

type Store struct {
	DB *database.Database
}

func New(config *configs.Config) *Store {
	db := NewDBClients(config.DSN)
	tempDB = db
	return &Store{DB: db}
}

func (s *Store) Shutdown() error {
	// gorm support auto shutdown
	return nil
}
