package model

import (
	"tirelease/internal/entity"
)

var _ IStateContext = (*PickTriageStateContext)(nil)
var _ IState[*PickTriageStateContext] = (*PickTriageState)(nil)

type PickTriageStateContext struct {
	VersionTriageID int64
	State           *PickTriageState
	Issue           *entity.Issue
	Version         *ReleaseVersion
	Prs             []entity.PullRequest
}

func (context *PickTriageStateContext) Trans(toState StateText) (bool, error) {
	isSuccess, err := context.State.Dispatch(context.GetStateText(), toState, context)
	if err != nil {
		return false, err
	}

	if context.Version.IsFrozen() && toState == ParseFromEntityPickTriage(entity.VersionTriageResultAccept) {
		toState = ParseFromEntityPickTriage(entity.VersionTriageResultAcceptFrozen)
	}
	context.State.setStateText(toState)
	return isSuccess, nil
}

func (context *PickTriageStateContext) GetStateText() StateText {
	return context.State.getStateText()
}

func (context *PickTriageStateContext) IsAccept() bool {
	return context.State.StateText == ParseFromEntityPickTriage(entity.VersionTriageResultAccept)
}

func NewPickTriageStateContext(stateText StateText, issue *entity.Issue,
	version *ReleaseVersion, prs []entity.PullRequest) (*PickTriageStateContext, error) {

	context := &PickTriageStateContext{}

	state, err := NewPickTriageState(stateText)
	if err != nil {
		return nil, err
	}
	context.State = state
	context.Issue = issue
	context.Version = version
	context.Prs = prs

	return context, nil
}

// Make the State struct private to force the only entrance be NewState func.
type PickTriageState struct {
	State[*PickTriageStateContext]
	StateText StateText
	transMap  TransitionMap[*PickTriageStateContext]
}

func NewPickTriageState(stateText StateText) (*PickTriageState, error) {
	state := &PickTriageState{
		StateText: stateText,
	}
	state.IState = interface{}(state).(IState[*PickTriageStateContext])
	state.init()

	return state, nil
}

func (state *PickTriageState) onLeave(context *PickTriageStateContext) (bool, error) {
	if state.StateText == ParseFromEntityPickTriage(entity.VersionTriageResultAccept) {
		for _, pr := range context.Prs {
			err := ChangePrApprovedLabel(pr, false, false)
			if err != nil {
				return false, nil
			}
		}
	}

	return true, nil
}

func (state *PickTriageState) getStateText() StateText {
	return state.StateText
}

func (state *PickTriageState) setStateText(stateText StateText) {
	state.StateText = stateText
}

func (state *PickTriageState) getTransition(meta StateTransitionMeta) IStateTransition[*PickTriageStateContext] {
	if state.transMap == nil {
		state.transMap = PickTriageTransMap
	}

	for k, v := range state.transMap {
		if k == meta {
			return v
		}
		if k.FromState == EmptyStateText() && k.ToState == meta.ToState {
			return v
		}
		if k.FromState == meta.FromState && k.ToState == EmptyStateText() {
			return v
		}
	}

	return PickTriageDefault{}
}

func (state *PickTriageState) init() error {
	if len(state.transMap) > 0 {
		return nil
	}
	state.transMap = PickTriageTransMap
	return nil
}
