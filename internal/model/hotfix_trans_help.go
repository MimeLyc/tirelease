package model

import (
	"fmt"
	"tirelease/commons/feishu"
	"tirelease/internal/constants"
	"tirelease/internal/entity"
	"tirelease/internal/service/notify"

	. "tirelease/commons/log"
)

func sendPendingApproveMessage(context hotfixStateContext) error {
	name := context.Hotfix.Name
	customer := "todo customer"
	repos := "todo repo"
	issues := "todo issue"
	prs := "todo pr"

	msg := getHotfixPendingApprovalMsg(name, customer, repos, issues, prs)
	header := "New Hotfix Creation Apply! 🔔"
	severity := constants.NotifySeverityInfo
	actions := []notify.Input{
		{
			Text: "Approve",
			Type: feishu.InteractiveTypePrimary,
			Value: map[string]interface{}{
				"register_object":    constants.EventRegisterObjectHotfix,
				"register_object_id": name,
				"register_action":    constants.EventRegisterApprove,
			},
		},
		{
			Text: "Deny",
			Type: feishu.InteractiveTypeDanger,
			Value: map[string]interface{}{
				"register_object":    constants.EventRegisterObjectHotfix,
				"register_object_id": name,
				"register_action":    constants.EventRegisterDeny,
			},
		},
	}

	notifyContent := notify.NotifyContent{
		Header: header,
		Blocks: []notify.Block{
			{
				Text:   msg,
				Inputs: actions,
			},
		},
		Severity: severity,
	}

	isActive := true
	eventRegistryCmd := EventRegistryCmd{
		Options: &entity.EventRegistryOptions{
			IsActive:    &isActive,
			EventObject: constants.EventRegisterObjectHotfix,
			EventAction: constants.EventRegisterPendingApproval,
			NotifyType:  constants.NotifyTriggerAction,
			ListOption:  entity.ListOption{},
		},
	}

	registries, err := eventRegistryCmd.BuildArray()
	if err != nil {
		return err
	}

	for _, registry := range registries {
		err := registry.Notify(notifyContent)
		if err != nil {
			Log.Errorf(err, "Notify error: %v", registry)
		}
	}

	return nil
}

func getHotfixPendingApprovalMsg(hotfixName, customer, repos, issues, pullrequest string) string {
	return fmt.Sprintf(
		constants.HotfixPendingApprovalMsg,
		hotfixName, customer, repos, issues, pullrequest,
	)
}

func sendHotfixApproveMessage(context hotfixStateContext) error {
	name := context.Hotfix.Name

	msg := fmt.Sprintf(constants.HotfixApproveMsg, name)
	header := "Congrats! 🔔"
	severity := constants.NotifySeverityInfo

	notifyContent := notify.NotifyContent{
		Header: header,
		Blocks: []notify.Block{
			{
				Text: msg,
				Links: []notify.Link{
					{
						Href: fmt.Sprintf("https://tirelease.pingcap.net/home/hotfix/%s", name),
						Text: name,
					}},
			},
		},
		Severity: severity,
	}

	isActive := true
	eventRegistryCmd := EventRegistryCmd{
		Options: &entity.EventRegistryOptions{
			IsActive:    &isActive,
			EventObject: constants.EventRegisterObjectHotfix,
			EventSpec:   name,
			EventAction: constants.EventRegisterApprove,
			NotifyType:  constants.NotifyTriggerAction,
		},
	}

	registries, err := eventRegistryCmd.BuildArray()
	if err != nil {
		return err
	}

	for _, registry := range registries {
		err := registry.Notify(notifyContent)
		if err != nil {
			Log.Errorf(err, "Notify error: %v", registry)
		}
	}

	return nil
}

func sendHotfixDenyMessage(context hotfixStateContext) error {
	name := context.Hotfix.Name

	msg := fmt.Sprintf(constants.HotfixDenyMsg, name)
	header := "Sorry! 🔔"
	severity := constants.NotifySeverityInfo

	notifyContent := notify.NotifyContent{
		Header: header,
		Blocks: []notify.Block{
			{
				Text: msg,
			},
		},
		Severity: severity,
	}

	isActive := true
	eventRegistryCmd := EventRegistryCmd{
		Options: &entity.EventRegistryOptions{
			IsActive:    &isActive,
			EventObject: constants.EventRegisterObjectHotfix,
			EventSpec:   name,
			EventAction: constants.EventRegisterDeny,
			NotifyType:  constants.NotifyTriggerAction,
		},
	}

	registries, err := eventRegistryCmd.BuildArray()
	if err != nil {
		return err
	}

	for _, registry := range registries {
		err := registry.Notify(notifyContent)
		if err != nil {
			Log.Errorf(err, "Notify error: %v", registry)
		}
	}

	return nil
}
