package model

import (
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

type HotfixCmd struct {
	HotfixOptions *entity.HotfixOptions
}

func (cmd HotfixCmd) Build() (*Hotfix, error) {
	hotfixes, err := repository.SelectFirstHotfixes(cmd.HotfixOptions)
	if err != nil {
		return nil, err
	}

	return &Hotfix{Hotfix: *hotfixes}, nil
}

func (cmd HotfixCmd) BuildArray() ([]Hotfix, error) {
	hotfixes, err := repository.SelectHotfixes(cmd.HotfixOptions)
	if err != nil {
		return nil, err
	}

	var models []Hotfix
	for _, hotfix := range *hotfixes {
		models = append(models, Hotfix{Hotfix: hotfix})
	}

	return models, nil
}

func (cmd HotfixCmd) Save(hotfix Hotfix) error {
	return repository.CreateOrUpdateHotfix(&hotfix.Hotfix)
}
