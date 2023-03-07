package model

import (
	"tirelease/internal/constants"
	"tirelease/internal/entity"

	. "tirelease/commons/log"
)

type Hotfix struct {
	entity.Hotfix
	Creator *User `json:"creator,omitempty"`
	HotfixArtifact
	ReleaseInfos []HotfixReleaseInfo `json:"release_infos,omitempty"`
}

type HotfixArtifact struct {
	ArtifactArchs    []string `json:"artifact_archs,omitempty"`
	ArtifactEditions []string `json:"artifact_editions,omitempty"`
	ArtifactTypes    []string `json:"artifact_types,omitempty"`
}

func (h *Hotfix) ChangeStatus(context hotfixStateContext) error {
	toStateText := StateText(context.ToState)

	_, err := context.Trans(toStateText)
	return err
}

// Call TiBuild to build related hotfix
func (h *Hotfix) Build() error {
	return nil
}

// DeActive all event registry related to pending_approval action of this hotfix.
func (h *Hotfix) TurnoffPendingApprovalNotify() error {
	isActive := true
	eventRegistryCmd := EventRegistryCmd{
		Options: &entity.EventRegistryOptions{
			IsActive:     &isActive,
			EventObject:  constants.EventRegisterObjectHotfix,
			EventSpec:    h.Name,
			EventActions: []constants.EventRegisterAction{constants.EventRegisterApprove, constants.EventRegisterDeny},
			NotifyType:   constants.NotifyTriggerAction,
		},
	}

	registries, err := eventRegistryCmd.BuildArray()
	if err != nil {
		return err
	}

	for _, registry := range registries {
		err := registry.DeActive()
		if err != nil {
			Log.Errorf(err, "Turnoff event registry notify error: %v", registry)
		}
	}

	return nil
}

// Register event registry for pending_approval action of this hotfix.
func (h *Hotfix) RegisterPendingApproval(user RegisterUser) error {
	register := EventRegistry{
		User: user,
		Event: RegisterEvent{
			Object: constants.EventRegisterObjectHotfix,
			Spec:   h.Name,
			Action: constants.EventRegisterApprove,
		},
		NotifyTrigger: NotifyTrigger{
			Type: constants.NotifyTriggerAction,
		},
		IsActive: true,
	}

	err := EventRegistryCmd{}.Save(register)

	if err != nil {
		return err
	}
	// Also register for denying action
	register.Event.Action = constants.EventRegisterDeny
	err = EventRegistryCmd{}.Save(register)
	if err != nil {
		return err
	}
	return nil
}
