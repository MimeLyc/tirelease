package model

import (
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

type EventRegistryRepo struct {
}

func (repo EventRegistryRepo) Save(registry EventRegistry) error {
	entity := mapEventRegistryToEntity(registry)
	return repository.CreateOrUpdateEventRegistry(&entity)
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
