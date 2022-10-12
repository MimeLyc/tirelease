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

func NewBlockTriageStateContext(stateText StateText,
	issue *entity.Issue, version *ReleaseVersion,
	historicalTriages *[]entity.VersionTriage) (*BlockTriageStateContext, error) {
	context := &BlockTriageStateContext{}

	state, err := NewBlockTriageState(stateText, issue, version, historicalTriages)
	if err != nil {
		return nil, err
	}

	context.State = state
	context.Issue = issue
	context.Version = version

	return context, nil
}

func (context *BlockTriageStateContext) Trans(toState StateText) (bool, error) {
	isSuccess, err := context.State.Dispatch(context.State.getStateText(), toState, context)
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

func NewBlockTriageState(stateText StateText,
	issue *entity.Issue, version *ReleaseVersion,
	historicalTriages *[]entity.VersionTriage) (*BlockTriageState, error) {
	if stateText == EmptyStateText() {
		stateText, _ = getBlockDefaultValue(issue, version, historicalTriages)
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

func getBlockDefaultValue(issue *entity.Issue,
	version *ReleaseVersion, historicalTriages *[]entity.VersionTriage) (StateText, error) {
	// default value of block triage status
	if issue.SeverityLabel == git.SeverityCriticalLabel {
		return ParseFromEntityBlockTriage(entity.BlockVersionReleaseResultBlock), nil
	}

	if len(*historicalTriages) > 0 {
		for _, triage := range *historicalTriages {
			if triage.VersionName < version.Name && triage.TriageResult == entity.VersionTriageResultReleased {
				return ParseFromEntityBlockTriage(entity.BlockVersionReleaseResultBlock), nil
			}
		}
	}

	return EmptyStateText(), nil
}
