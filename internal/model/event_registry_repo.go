package model

import (
	"tirelease/internal/entity"
	"tirelease/internal/store"
)

type EventRegistryCmd struct {
	Options *entity.EventRegistryOptions
}

func (repo EventRegistryCmd) BuildArray() ([]EventRegistry, error) {
	entities, err := store.SelectEventRegistries(repo.Options)
	if err != nil {
		return nil, err
	}

	var models []EventRegistry
	for _, entity := range *entities {
		models = append(
			models,
			EventRegistry{
				User: RegisterUser{
					Platform: entity.UserPlatform,
					Type:     entity.UserType,
					ID:       entity.UserID,
				},
				Event: RegisterEvent{
					Object: entity.EventObject,
					Spec:   entity.EventSpec,
					Action: entity.EventAction,
				},
				NotifyTrigger: NotifyTrigger{
					Type:   entity.NotifyType,
					Config: entity.NotifyConfig,
				},
				IsActive: entity.IsActive,
			},
		)
	}

	return models, nil
}

func (repo EventRegistryCmd) Save(registry EventRegistry) error {
	entity := mapEventRegistryToEntity(registry)
	return store.CreateOrUpdateEventRegistry(&entity)
}

func mapEventRegistryToEntity(registry EventRegistry) entity.EventRegistry {
	return entity.EventRegistry{
		ID:           0,
		IsDeleted:    false,
		IsActive:     registry.IsActive,
		UserPlatform: registry.User.Platform,
		UserType:     registry.User.Type,
		UserID:       registry.User.ID,
		EventObject:  registry.Event.Object,
		EventSpec:    registry.Event.Spec,
		EventAction:  registry.Event.Action,
		NotifyType:   registry.NotifyTrigger.Type,
		NotifyConfig: registry.NotifyTrigger.Config,
	}
}
