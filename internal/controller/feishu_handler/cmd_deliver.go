package feishu_handler

import (
	"fmt"
	"tirelease/internal/constants"
	"tirelease/internal/model"
	"tirelease/internal/service"
	"tirelease/internal/service/notify"
)

// Chatops command pattern:
//
//	object(issue, pr, version) + cmd(aprove, watch...) + flags
func deliverCmd(receive MsgReceiveV1) error {
	content, err := NewContent(receive.Event.Message.Content)

	// TODO validate content
	// content.validate()...

	switch content.cmd {
	case "watch":
		err = registerEvent(receive, content)
	}

	replyMessage(receive, content, err)
	return err
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
		IsActive:      true,
	}

	return service.RegisterEvent(registry)
}

func replyMessage(receive MsgReceiveV1, content content, err error) error {
	target := content.target
	cmd := content.cmd
	severity := constants.NotifySeverityInfo
	header := "Congrats🥳 !"

	msg := fmt.Sprintf("You have successfully **%sed** target **%s**!", cmd, target)
	if err != nil {
		header = "Sorry🙏 !"
		severity = constants.NotifySeverityAlarm
		msg = fmt.Sprintf("You failed to **%s** the **%s**!\n"+
			"The error msg is:\n"+
			"<font color='green'>%s</font>\n"+
			"You can ask the developer of TiRelease for help.\n",
			cmd, target, err.Error())
	}

	block := notify.Block{
		Text: msg,
	}

	notifyContent := notify.NotifyContent{
		Header:   header,
		Severity: severity,
		Blocks:   []notify.Block{block},
	}

	return notify.ReplyFeishuByMessageId(receive.Event.Message.MessageID, notifyContent)

}
