package service

import (
	"tirelease/internal/dto"
	"tirelease/internal/entity"
	"tirelease/internal/model"
	"tirelease/internal/store"
)

// Save function save the **whole** request.Hotfix to database.
func SaveHotfix(request dto.HotfixSaveRequest) (*model.Hotfix, error) {
	newHotfix, err := buildHotfixFromRequest(request)
	if err != nil {
		return nil, err
	}

	hotfixCmd := model.HotfixCmd{
		HotfixOptions: &entity.HotfixOptions{
			Name: newHotfix.Name,
		},
	}

	old, err := hotfixCmd.Build()
	if err != nil {
		if err, ok := err.(store.DataNotFoundError); !ok {
			return nil, err
		}
	}

	// revert status and make the changing
	if old == nil {
		// Dumb status for trigger creating new hotfix events.
		newHotfix.Status = entity.HotfixStatusInit
	} else {
		newHotfix.Status = old.Status
	}
	//	Change hotfix status
	context := model.NewHotfixStateContext(&newHotfix)
	context.ToState = request.Status
	context.OperatorEmail = request.OperatorEmail
	err = newHotfix.ChangeStatus(context)

	if err != nil {
		return nil, err
	}

	return &newHotfix, hotfixCmd.Save(newHotfix)
}

func FindHotfixes(options entity.HotfixOptions) ([]model.Hotfix, error) {
	hotfixCmd := model.HotfixCmd{
		HotfixOptions: &options,
	}

	return hotfixCmd.BuildArray()
}

func FindHotfixByName(name string) (*model.Hotfix, error) {
	hotfixCmd := model.HotfixCmd{
		HotfixOptions: &entity.HotfixOptions{
			Name: name,
		},
	}

	return hotfixCmd.Build()
}

// ApproveHotfix function approve pending_aproval hotfix to upcoming.
func ApproveHotfix(name string) error {
	hotfixCmd := model.HotfixCmd{
		HotfixOptions: &entity.HotfixOptions{
			Name: name,
		},
	}

	hotfix, err := hotfixCmd.Build()
	if hotfix == nil || err != nil {
		return err
	}

	if hotfix.Status != entity.HotfixStatusUpcoming {
		context := model.NewHotfixStateContext(hotfix)
		context.ToState = entity.HotfixStatusUpcoming
		err := hotfix.ChangeStatus(context)

		if err != nil {
			return err
		}
	}

	return hotfixCmd.Save(*hotfix)

}

// DenyHotfix function deny pending_aproval hotfix to denied.
func DenyHotfix(name string) error {
	hotfixCmd := model.HotfixCmd{
		HotfixOptions: &entity.HotfixOptions{
			Name: name,
		},
	}

	hotfix, err := hotfixCmd.Build()
	if hotfix == nil || err != nil {
		return err
	}

	if hotfix.Status != entity.HotfixStatusDenied {
		context := model.NewHotfixStateContext(hotfix)
		context.ToState = entity.HotfixStatusDenied
		err := hotfix.ChangeStatus(context)

		if err != nil {
			return err
		}
	}

	return hotfixCmd.Save(*hotfix)
}
