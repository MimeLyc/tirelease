package model

import (
	. "tirelease/commons/log"
	"tirelease/internal/entity"
)

// Status to Won't Fix
type PickTriage2WontFix struct {
}

func (trans PickTriage2WontFix) FitConstraints(context *pickTriageStateContext) (bool, error) {
	return true, nil
}

func (trans PickTriage2WontFix) Effect(context *pickTriageStateContext) (bool, error) {
	prs := context.Prs
	// Derectly close all PRs related to this issue.
	for _, pr := range prs {
		err := pr.Close()

		if err != nil {
			Log.Errorf(err, "close pr %d failed, err: %v", pr.ID, err)
		}
	}

	return true, nil
}

// Status to Won't Fix End

// Status to Approved
type PickTriage2Accept struct {
}

func (trans PickTriage2Accept) FitConstraints(context *pickTriageStateContext) (bool, error) {
	return true, nil
}

func (trans PickTriage2Accept) Effect(context *pickTriageStateContext) (bool, error) {
	isFrozen := context.Version.IsFrozen()

	for _, pr := range context.Prs {
		pr := pr
		if isFrozen {
			if err := pr.UnApprove(); err != nil {
				return false, err
			}

		} else {
			if err := pr.Approve(); err != nil {
				return false, err
			}
		}
	}

	return true, nil
}

// Status to Approved
type PickTriage2AcceptFrozen struct {
}

func (trans PickTriage2AcceptFrozen) FitConstraints(context *pickTriageStateContext) (bool, error) {
	return true, nil
}

func (trans PickTriage2AcceptFrozen) Effect(context *pickTriageStateContext) (bool, error) {
	for _, pr := range context.Prs {
		pr := pr
		err := pr.UnApprove()
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

// Status to Approved End

type PickTriageDefault struct {
}

func (trans PickTriageDefault) FitConstraints(context *pickTriageStateContext) (bool, error) {
	return true, nil
}

func (trans PickTriageDefault) Effect(context *pickTriageStateContext) (bool, error) {
	return true, nil
}

var PickTriageTransMap = make(TransitionMap[*pickTriageStateContext])

// Orders matters, the trans with from and to should be added firstly.
func init() {
	if len(PickTriageTransMap) > 0 {
		return
	}
	PickTriageTransMap[StateTransitionMeta{
		FromState: EmptyStateText(),
		ToState:   StateText(entity.VersionTriageResultWontFix),
	}] = PickTriage2WontFix{}

	PickTriageTransMap[StateTransitionMeta{
		FromState: EmptyStateText(),
		ToState:   StateText(entity.VersionTriageResultAccept),
	}] = PickTriage2Accept{}

	PickTriageTransMap[StateTransitionMeta{
		FromState: EmptyStateText(),
		ToState:   StateText(entity.VersionTriageResultAcceptFrozen),
	}] = PickTriage2AcceptFrozen{}

}
