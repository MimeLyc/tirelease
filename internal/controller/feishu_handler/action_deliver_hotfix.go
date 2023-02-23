package feishu_handler

import (
	"fmt"
	"tirelease/internal/constants"
	"tirelease/internal/service"
)

func deliverHotfixAction(request ActionRequest) error {
	hotfixName := request.ObjectID
	action := request.Action

	if action == constants.EventRegisterApprove {
		return service.ApproveHotfix(hotfixName)
	} else if action == constants.EventRegisterDeny {
		return service.DenyHotfix(hotfixName)
	}

	return fmt.Errorf("unknown action: %s", action)
}
