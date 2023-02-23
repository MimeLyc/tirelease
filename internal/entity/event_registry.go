package entity

import (
	"tirelease/internal/constants"
)

type EventRegistry struct {
	// DataBase Column
	ID int64 `json:"id,omitempty"`

	// Soft delete by tag
	IsDeleted bool `json:"is_deleted,omitempty"`

	IsActive bool `json:"is_active,omitempty" default:"true"`

	// Register user info
	UserPlatform constants.EventRegisterPlatform `json:"user_platform,omitempty"`
	UserType     constants.EventRegisterUserType `json:"user_type,omitempty"`
	UserID       string                          `json:"user_id,omitempty"`

	// Register event Info
	EventObject constants.EventRegisterObject `json:"event_object,omitempty"`
	EventSpec   string                        `json:"event_spec,omitempty"`
	EventAction constants.EventRegisterAction `json:"event_action,omitempty"`

	// NotifyTrigger
	NotifyType   constants.NotifyTriggerType `json:"notify_type,omitempty"`
	NotifyConfig string                      `json:"notify_config,omitempty"`
}

func (EventRegistry) TableName() string {
	return "event_registry"
}

type EventRegistryOptions struct {
	IsActive *bool `json:"is_active,omitempty"`

	// Register user info
	UserPlatform constants.EventRegisterPlatform `json:"user_platform,omitempty"`
	UserType     constants.EventRegisterUserType `json:"user_type,omitempty"`
	UserID       string                          `json:"user_id,omitempty"`

	// Register event Info
	EventObject constants.EventRegisterObject `json:"event_object,omitempty"`
	EventSpec   string                        `json:"event_spec,omitempty"`
	EventAction constants.EventRegisterAction `json:"event_action,omitempty"`

	// NotifyTrigger
	NotifyType   constants.NotifyTriggerType `json:"notify_type,omitempty"`
	NotifyConfig string                      `json:"notify_config,omitempty"`

	EventActions []constants.EventRegisterAction `json:"event_actions,omitempty"`

	ListOption
}
