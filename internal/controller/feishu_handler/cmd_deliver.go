package feishu_handler

import (
	"tirelease/internal/constants"
	"tirelease/internal/model"
	"tirelease/internal/service"
)

// Chatops command pattern:
//
//	object(issue, pr, version) + cmd(aprove, watch...) + flags
func deliverCmd(receive MsgReceiveV1) error {
	content, err := NewContent(receive.Event.Message.Content)
	if err != nil {
		return err
	}

	// TODO validate content
	// content.validate()...

	switch content.cmd {
	case "watch":
		return registerEvent(receive, content)
	}

	// msgId := receive.Event.Message.MessageID
	// todo reply unfind mesg
	return nil
}

func registerEvent(receive MsgReceiveV1, content content) error {
	user := model.RegisterUser{
		Platform: constants.EventRegisterPlatformFeishu,
		Type:     constants.EventRegisterUserType(receive.Event.Message.ChatType),
		ID:       receive.Event.Message.ChatID,
	}

	target := content.target
	action := content.extractByKey("action")
	spec := content.extractSpec()
	event := model.RegisterEvent{
		Object: constants.EventRegisterObject(target),
		Spec:   spec,
		Action: constants.EventRegisterAction(action),
	}

	trigger := model.NotifyTrigger{
		Type: constants.NotifyTriggerAction,
	}

	registry := model.EventRegistry{
		User:          user,
		Event:         event,
		NotifyTrigger: trigger,
	}

	return service.RegisterEvent(registry)
}
