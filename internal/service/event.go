package service

import "tirelease/internal/model"

func RegisterEvent(event model.EventRegistry) error {
	repo := model.EventRegistryCmd{}
	err := repo.Save(event)

	// todo ,send msg to register user
	if err != nil {
		return err
	}

	return nil
}
