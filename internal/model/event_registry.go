package model

import (
	"tirelease/internal/constants"
)

type EventRegistry struct {
	User  RegisterUser
	Event RegisterEvent
	NotifyTrigger
	IsActive bool
}

type RegisterUser struct {
	Platform constants.EventRegisterPlatform
	Type     constants.EventRegisterUserType
	ID       string
}

type RegisterEvent struct {
	Object constants.EventRegisterObject
	Spec   string
	Action constants.EventRegisterAction
}

type NotifyTrigger struct {
	Type   constants.NotifyTriggerType
	Config string
}
