package model

import (
	"tirelease/commons/git"
	"tirelease/internal/entity"
)

var _ IStateContext = (*BlockTriageStateContext)(nil)
var _ IState[*BlockTriageStateContext] = (*BlockTriageState)(nil)

type BlockTriageStateContext struct {
	State   *BlockTriageState
	Issue   *entity.Issue
	Version *ReleaseVersion
}

func NewBlockTriageStateContext(stateText StateText, issue *entity.Issue, version *ReleaseVersion) (*BlockTriageStateContext, error) {
	context := &BlockTriageStateContext{}

	state, err := NewBlockTriageState(stateText, issue, version)
	if err != nil {
		return nil, err
	}

	context.State = state
	context.Issue = issue
	context.Version = version

	return context, nil
}

func (context *BlockTriageStateContext) Trans(toState StateText) (bool, error) {
	isSuccess, err := context.State.Dispatch(toState, context)
	if err != nil {
		return false, err
	}

	context.State.setStateText(toState)
	return isSuccess, nil
}

// State
type BlockTriageState struct {
	State[*BlockTriageStateContext]
	StateText StateText
	transMap  TransitionMap[*BlockTriageStateContext]
}

func NewBlockTriageState(stateText StateText, issue *entity.Issue, version *ReleaseVersion) (*BlockTriageState, error) {
	if stateText == EmptyStateText() {
		// default value of block triage status
		if issue.SeverityLabel == git.SeverityCriticalLabel {
			stateText = ParseFromEntityBlockTriage(entity.BlockVersionReleaseResultBlock)
		}
	}

	state := &BlockTriageState{
		StateText: stateText,
	}
	state.IState = interface{}(state).(IState[*BlockTriageStateContext])
	// state.init()

	return state, nil
}

func (state *BlockTriageState) getStateText() StateText {
	return state.StateText
}

func (state *BlockTriageState) setStateText(stateText StateText) {
	state.StateText = stateText
}
