package model

// State is the template of state machine.
// Constructed by template model, see https://refactoring.guru/design-patterns/template-method/go/example
type State[T IStateContext] struct {
	IState[T]
}

func (state State[T]) onTran(trans IStateTransition[T], context T) (bool, error) {
	fitConstraints, err := trans.FitConstraints(context)

	if err != nil {
		return false, nil
	}

	if fitConstraints {
		isTransOk, err := trans.Effect(context)
		if err != nil {
			return false, nil
		}

		return isTransOk, nil

	} else {
		return false, nil
	}
}

func (state State[T]) init() error {
	return nil
}

func (state State[T]) getStateText() StateText {
	return state.IState.getStateText()
}

func (state State[T]) onLeave(context T) (bool, error) {
	return true, nil
}

func (state State[T]) getTransition(meta StateTransitionMeta) IStateTransition[T] {
	return nil
}

func (state State[T]) Dispatch(fromState, toState StateText, context T) (bool, error) {
	if toState == fromState {
		return false, nil
	}

	isLeaveOK, err := state.IState.onLeave(context)
	if err != nil {
		return false, err
	}
	if !isLeaveOK {
		return false, nil
	}

	transition := state.IState.getTransition(StateTransitionMeta{state.getStateText(), toState})
	if transition == nil {
		return true, nil
	}

	isTransOK, err := state.IState.onTran(transition, context)
	if err != nil {
		return false, err
	}
	if !isTransOK {
		return false, nil
	}

	return true, nil
}
