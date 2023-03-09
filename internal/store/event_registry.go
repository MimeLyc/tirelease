package store

import (
	"fmt"
	"tirelease/commons/database"
	"tirelease/internal/entity"

	. "tirelease/commons/log"

	"github.com/pkg/errors"
	"gorm.io/gorm/clause"
)

func SelectEventRegistries(option *entity.EventRegistryOptions) (*[]entity.EventRegistry, error) {
	sql := "select * from event_registry where 1=1" + eventRegistryWhere(option) + option.GetOrderByString() + option.GetLimitString()
	// 查询
	var eventRegistries []entity.EventRegistry
	if err := database.DBConn.RawWrapper(sql, option).Find(&eventRegistries).Error; err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("select event regitstry: %+v failed", option))
	}

	return &eventRegistries, nil
}

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

func eventRegistryWhere(option *entity.EventRegistryOptions) string {
	sql := ""

	if option.IsActive != nil {
		sql += " and event_registry.is_active = @IsActive"
	}
	if option.UserPlatform != "" {
		sql += " and event_registry.user_platform = @UserPlatform"
	}
	if option.UserType != "" {
		sql += " and event_registry.user_type = @UserType"
	}
	if option.UserID != "" {
		sql += " and event_registry.user_id = @UserID"
	}
	if option.EventObject != "" {
		sql += " and event_registry.event_object = @EventObject"
	}
	if option.EventAction != "" {
		sql += " and event_registry.event_action = @EventAction"
	}
	if option.EventSpec != "" {
		sql += " and event_registry.event_spec = @EventSpec"
	}
	if option.NotifyType != "" {
		sql += " and event_registry.notify_type = @NotifyType"
	}
	if option.NotifyConfig != "" {
		sql += " and event_registry.notify_config = @NotifyConfig"
	}
	if len(option.EventActions) != 0 {
		sql += " and event_registry.event_action in @EventActions"
	}

	return sql
}
