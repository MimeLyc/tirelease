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
	msg := getHotfixPendingApprovalMsg(*context.Hotfix)

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
				Text: msg,
				Links: []notify.Link{
					{Text: "TiReleaes Hotfix Page",
						Href: fmt.Sprintf("%s/home/hotfix/%s", constants.TiReleaseUrl, name),
					},
				},
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

func getHotfixPendingApprovalMsg(hotfix Hotfix) string {
	repos := ""
	for _, release := range hotfix.ReleaseInfos {
		repos += release.RepoFullName
	}

	return fmt.Sprintf(
		constants.HotfixPendingApprovalMsg,
		hotfix.Name,
		hotfix.Customer,
		hotfix.Creator.Name,
		hotfix.OncallUrl,
		fmt.Sprintf("%s-%s", hotfix.OncallPrefix, hotfix.OncallID),
		hotfix.BaseVersionName,
		repos,
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
						Href: fmt.Sprintf("%s/home/hotfix/%s", constants.TiReleaseUrl, name),
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
