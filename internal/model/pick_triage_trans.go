package model

import (
	"tirelease/internal/entity"
)

// Status to Won't Fix
type PickTriage2WontFix struct {
}

func (trans PickTriage2WontFix) FitConstraints(context *PickTriageStateContext) (bool, error) {
	return true, nil
}

func (trans PickTriage2WontFix) Effect(context *PickTriageStateContext) (bool, error) {
	prs := context.Prs
	version := context.Version
	issue := context.Issue
	err := closeWontfixPrs(prs, issue.IssueID, version.Name)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Status to Won't Fix End

// Status to Approved
type PickTriage2Accept struct {
}

func (trans PickTriage2Accept) FitConstraints(context *PickTriageStateContext) (bool, error) {
	return true, nil
}

func (trans PickTriage2Accept) Effect(context *PickTriageStateContext) (bool, error) {
	isFrozen := context.Version.IsFrozen()

	for _, pr := range context.Prs {
		pr := pr
		err := ChangePrApprovedLabel(pr, isFrozen, true)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

// Status to Approved
type PickTriage2AcceptFrozen struct {
}

func (trans PickTriage2AcceptFrozen) FitConstraints(context *PickTriageStateContext) (bool, error) {
	return true, nil
}

func (trans PickTriage2AcceptFrozen) Effect(context *PickTriageStateContext) (bool, error) {
	for _, pr := range context.Prs {
		pr := pr
		err := ChangePrApprovedLabel(pr, true, true)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

// Status to Approved End

type PickTriageDefault struct {
}

func (trans PickTriageDefault) FitConstraints(context *PickTriageStateContext) (bool, error) {
	return true, nil
}

func (trans PickTriageDefault) Effect(context *PickTriageStateContext) (bool, error) {
	return true, nil
}

var PickTriageTransMap = make(TransitionMap[*PickTriageStateContext])

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
