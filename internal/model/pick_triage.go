package model

import (
	"tirelease/internal/entity"
)

var _ IStateContext = (*pickTriageStateContext)(nil)
var _ IState[*pickTriageStateContext] = (*PickTriageState)(nil)

type pickTriageStateContext struct {
	VersionTriageID int64
	State           *PickTriageState
	Issue           *entity.Issue
	Version         *ReleaseVersion
	Prs             []PullRequest
	IsForce         bool
}

func (context *pickTriageStateContext) Trans(toState StateText) (bool, error) {
	isForce := context.IsForce
	isFrozen := context.Version.IsFrozen()
	isApprove := toState == ParseFromEntityPickTriage(entity.VersionTriageResultAccept)
	if isApprove {
		if !isForce && isFrozen {
			toState = ParseFromEntityPickTriage(entity.VersionTriageResultAcceptFrozen)
		}
	}
	isSuccess, err := context.State.Dispatch(context.GetStateText(), toState, context)
	if err != nil {
		return false, err
	}

	context.State.setStateText(toState)
	return isSuccess, nil
}

func (context *pickTriageStateContext) GetStateText() StateText {
	return context.State.getStateText()
}

func (context *pickTriageStateContext) IsAccept() bool {
	return context.State.StateText == ParseFromEntityPickTriage(entity.VersionTriageResultAccept)
}

func NewPickTriageStateContext(stateText StateText, issue *entity.Issue,
	version *ReleaseVersion, prs []PullRequest) (*pickTriageStateContext, error) {

	context := &pickTriageStateContext{}

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
	State[*pickTriageStateContext]
	StateText StateText
	transMap  TransitionMap[*pickTriageStateContext]
}

func NewPickTriageState(stateText StateText) (*PickTriageState, error) {
	state := &PickTriageState{
		StateText: stateText,
	}
	state.IState = interface{}(state).(IState[*pickTriageStateContext])
	state.init()

	return state, nil
}

func (state *PickTriageState) onLeave(context *pickTriageStateContext) (bool, error) {
	if state.StateText == ParseFromEntityPickTriage(entity.VersionTriageResultAccept) {
		for _, pr := range context.Prs {
			err := pr.UnApprove()
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

func (state *PickTriageState) getTransition(meta StateTransitionMeta) IStateTransition[*pickTriageStateContext] {
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
