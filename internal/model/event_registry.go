package model

import (
	"fmt"
	"tirelease/internal/constants"
	"tirelease/internal/service/notify"
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

func (e EventRegistry) Notify(content notify.NotifyContent) error {
	notifyType := e.User.Type
	notifyID := e.User.ID
	notifyPlatform := e.User.Platform
	if notifyPlatform == constants.EventRegisterPlatformFeishu {
		switch notifyType {
		case constants.EventRegisterGroup:
			return notify.SendFeishuFormattedByGroup(notifyID, content)
		case constants.EventRegisterP2P:
			return notify.SendFeishuFormattedByEmail(notifyID, content)
		}
	}
	return fmt.Errorf("Invalid notification config : %v", e)
}

func (e EventRegistry) DeActive() error {
	e.IsActive = false
	return EventRegistryCmd{}.Save(e)
}
