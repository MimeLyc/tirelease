package feishu_handler

import (
	"fmt"
	"tirelease/internal/constants"
)

func deliverAction(action ActionReceive) error {
	if err := action.Validate(); err != nil {
		return err
	}
	request := NewActionRequest(action)

	switch request.Object {
	case constants.EventRegisterObjectHotfix:
		return deliverHotfixAction(request)
	case constants.EventRegisterObjectIssue:
	}

	return fmt.Errorf("Unknown object: %v", request.Object)
}
