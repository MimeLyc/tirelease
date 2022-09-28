package model

import (
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

var _ IStateContext = (*VersionStateContext)(nil)
var _ IState[*VersionStateContext] = (*VersionState)(nil)

type VersionStateContext struct {
	Version *ReleaseVersion
	State   *VersionState
}

func (context *VersionStateContext) Trans(toState StateText) (bool, error) {
	fromState := context.GetStateText()
	context.Version.ReleaseVersion.Status = entity.ReleaseVersionStatus(toState)
	// Update version status first because the next version may rely on the real status.
	repository.UpdateReleaseVersion(context.Version.ReleaseVersion)
	context.State.setStateText(toState)

	isSuccess, err := context.State.Dispatch(fromState, toState, context)
	if err != nil {
		return false, err
	}

	return isSuccess, nil
}

func (context *VersionStateContext) GetStateText() StateText {
	return context.State.getStateText()
}

func NewVersionStateContext(version *ReleaseVersion) (*VersionStateContext, error) {
	context := &VersionStateContext{}

	state, err := NewVersionState(StateText(version.Status))
	if err != nil {
		return nil, err
	}
	context.State = state
	context.Version = version

	return context, nil
}

// Make the State struct private to force the only entrance be NewState func.
type VersionState struct {
	State[*VersionStateContext]
	StateText StateText
	transMap  TransitionMap[*VersionStateContext]
}

func NewVersionState(stateText StateText) (*VersionState, error) {
	state := &VersionState{
		StateText: stateText,
	}
	state.IState = interface{}(state).(IState[*VersionStateContext])
	state.init()

	return state, nil
}

func (state *VersionState) getStateText() StateText {
	return state.StateText
}

func (state *VersionState) setStateText(stateText StateText) {
	state.StateText = stateText
}

func (state *VersionState) getTransition(meta StateTransitionMeta) IStateTransition[*VersionStateContext] {
	if state.transMap == nil {
		state.transMap = VersionStateTransMap
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

	return nil
}

func (state *VersionState) init() error {
	if len(state.transMap) > 0 {
		return nil
	}
	state.transMap = VersionStateTransMap
	return nil
}
