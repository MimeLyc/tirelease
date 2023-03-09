package store

import (
	"tirelease/internal/entity"
	"tirelease/utils/configs"
	"tirelease/utils/database"
)

var tempDB *database.Database
var tempHrEmployeeDB *database.Database

func NewDBClients(dsn string) *database.Database {
	db := database.MustConnect(dsn)
	return db
}

type Store struct {
	DB *database.Database
}

func NewStore(config *configs.Config) *Store {
	db := NewDBClients(config.DSN)
	tempDB = db
	tempHrEmployeeDB = NewDBClients(config.EmployeeDSN)

	if config.RunAutoMigrate {
		db.AutoMigrate(
			&entity.SprintMeta{},
			&entity.PullRequest{},
			&entity.EventRegistry{},
			&entity.Hotfix{},
			&entity.HotfixReleaseInfo{},
		)
	}
	return &Store{DB: db}
}

func (s *Store) Shutdown() error {
	// gorm support auto shutdown
	return nil
}
