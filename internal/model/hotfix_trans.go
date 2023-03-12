package model

import (
	. "tirelease/commons/log"
	"tirelease/internal/constants"
	"tirelease/internal/entity"
)

var hotfixStateTransMap = make(TransitionMap[*hotfixStateContext])

func init() {
	if len(hotfixStateTransMap) > 0 {
		return
	}

	hotfixStateTransMap[StateTransitionMeta{
		FromState: EmptyStateText(),
		ToState:   StateText(entity.HotfixStatusPendingApproval),
	}] = hotfix2PendingApproval{}
	hotfixStateTransMap[StateTransitionMeta{
		FromState: EmptyStateText(),
		ToState:   StateText(entity.HotfixStatusDenied),
	}] = hotfix2Deny{}
	hotfixStateTransMap[StateTransitionMeta{
		FromState: EmptyStateText(),
		ToState:   StateText(entity.HotfixStatusUpcoming),
	}] = hotfix2Upcoming{}

}

type hotfix2Deny struct{}

func (trans hotfix2Deny) FitConstraints(context *hotfixStateContext) (bool, error) {
	return true, nil
}

func (trans hotfix2Deny) Effect(context *hotfixStateContext) (bool, error) {
	// Send pending approval notification
	err := sendHotfixDenyMessage(*context)
	if err != nil {
		return false, err
	}

	return true, context.Hotfix.TurnoffPendingApprovalNotify()
}

type hotfix2Upcoming struct{}

func (trans hotfix2Upcoming) FitConstraints(context *hotfixStateContext) (bool, error) {
	return true, nil
}

func (trans hotfix2Upcoming) Effect(context *hotfixStateContext) (bool, error) {
	// Checkout branch
	hotfix := context.Hotfix
	for _, release := range hotfix.ReleaseInfos {
		branch, err := release.FetchHotfixBranch(hotfix.BaseVersionName)
		if err != nil {
			Log.Errorf(err, "Fetch hotfix branch error: %v", release)
			continue
		}
		release.Branch = branch

		err = HotfixReleaseCmd{}.Save(release)
		if err != nil {
			return false, err
		}
	}

	// Send pending approval notification
	err := sendHotfixApproveMessage(*context)
	if err != nil {
		return false, err
	}

	return true, context.Hotfix.TurnoffPendingApprovalNotify()
}

type hotfix2PendingApproval struct{}

func (trans hotfix2PendingApproval) FitConstraints(context *hotfixStateContext) (bool, error) {
	return true, nil
}

func (trans hotfix2PendingApproval) Effect(context *hotfixStateContext) (bool, error) {
	// Send pending approval notification
	err := sendPendingApproveMessage(*context)
	if err != nil {
		return false, err
	}

	// Register event to event registry
	err = context.Hotfix.RegisterPendingApproval(
		RegisterUser{
			Platform: constants.EventRegisterPlatformFeishu,
			Type:     constants.EventRegisterP2P,
			ID:       context.OperatorEmail,
		},
	)
	if err != nil {
		return false, err
	}

	return true, nil
}
