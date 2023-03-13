package model

import (
	"tirelease/internal/entity"
)

// FSM state changing context for hotfix.
type hotfixStateContext struct {
	Hotfix        *Hotfix
	state         *hotfixState
	ToState       entity.HotfixStatus
	OperatorEmail string
}

func NewHotfixStateContext(hotfix *Hotfix) hotfixStateContext {
	context := hotfixStateContext{}
	state, err := NewHotfixState(StateText(hotfix.Status))
	if err != nil {
		return context
	}
	context.state = state
	context.Hotfix = hotfix

	return context
}

// Trans method changes the state of hotfix and dispatch the state transition event.
func (context *hotfixStateContext) Trans(toState StateText) (bool, error) {
	fromState := context.GetStateText()
	context.Hotfix.Status = entity.HotfixStatus(toState)
	context.state.setStateText(toState)

	isSuccess, err := context.state.Dispatch(fromState, toState, context)
	if err != nil {
		return false, err
	}

	return isSuccess, nil
}

func (context *hotfixStateContext) GetStateText() StateText {
	return context.state.getStateText()
}

func (context *hotfixStateContext) GetToStateText() (StateText, error) {
	hotfix := context.Hotfix
	fromStatus := hotfix.Status

	if fromStatus == entity.HotfixStatusQATesting {
		if hotfix.PassedQATest() {
			return StateText(entity.HotfixStatusReleased), nil
		} else if hotfix.PassQATest == entity.QATestResultFailed {
			return StateText(entity.HotfixStatusUpcoming), nil
		} else {
			return StateText(entity.HotfixStatusQATesting), nil
		}
	}

	// TODO: @mofeihu Add building logic and upcoming logic
	// if fromStatus == entity.HotfixStatusBuilding {
	// 	if building success, return qa testing
	// }

	if fromStatus == entity.HotfixStatusUpcoming {
		// If hotfix is ready for building, return Building
	}

	return StateText(context.ToState), nil
}

// hotfixState get the corresponding transition and dispatch the event.
// VO, unmutable
type hotfixState struct {
	State[*hotfixStateContext]
	StateText StateText
	transMap  TransitionMap[*hotfixStateContext]
}

func NewHotfixState(stateText StateText) (*hotfixState, error) {
	state := &hotfixState{
		StateText: stateText,
	}
	state.IState = interface{}(state).(IState[*hotfixStateContext])
	state.init()

	return state, nil
}

func (state *hotfixState) getStateText() StateText {
	return state.StateText
}

func (state *hotfixState) setStateText(stateText StateText) {
	state.StateText = stateText
}

func (state *hotfixState) getTransition(meta StateTransitionMeta) IStateTransition[*hotfixStateContext] {
	if state.transMap == nil {
		state.transMap = hotfixStateTransMap
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

func (state *hotfixState) init() error {
	if len(state.transMap) > 0 {
		return nil
	}
	state.transMap = hotfixStateTransMap
	return nil
}
