package model

import (
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

type HotfixCmd struct {
	HotfixOptions *entity.HotfixOptions
}

func (cmd HotfixCmd) Build() (*Hotfix, error) {
	hotfix, err := repository.SelectFirstHotfix(cmd.HotfixOptions)
	if err != nil {
		return nil, err
	}

	releaseInfos, err := HotfixReleaseCmd{
		HotfixReleaseInfoOptions: &entity.HotfixReleaseInfoOptions{
			HotfixName: hotfix.Name,
		},
	}.BuildArray()
	if err != nil {
		return nil, err
	}

	creator, err := UserCmd{}.BuildByEmail(hotfix.CreatorEmail)
	if err != nil {
		return nil, err
	}

	return &Hotfix{
		Hotfix:  *hotfix,
		Creator: creator,
		HotfixArtifact: HotfixArtifact{
			ArtifactArchs:    hotfix.UnserializeArtifactArchs(),
			ArtifactEditions: hotfix.UnserializeArtifactEditions(),
			ArtifactTypes:    hotfix.UnserializeArtifactTypes(),
		},
		ReleaseInfos: releaseInfos,
	}, nil
}

// TODO: Fill hotfix releaseInfos
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
