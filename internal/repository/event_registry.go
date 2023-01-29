package repository

import (
	"fmt"
	"tirelease/commons/database"
	"tirelease/internal/entity"

	. "tirelease/commons/log"

	"github.com/pkg/errors"
	"gorm.io/gorm/clause"
)

func CreateOrUpdateEventRegistry(registry *entity.EventRegistry) error {
	// 存储
	if err := database.DBConn.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&registry).Error; err != nil {
		Log.Errorf(err, "Create or update event_registry %v", registry)
		return errors.Wrap(err, fmt.Sprintf("create or update event_registry: %+v failed", registry))
	}

	Log.Infof("Create or update event_registry %v", registry)
	return nil
}
