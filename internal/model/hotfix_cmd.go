package model

import (
	"strings"
	"tirelease/internal/entity"
	"tirelease/internal/store"
)

type HotfixCmd struct {
	HotfixOptions *entity.HotfixOptions
}

func (cmd HotfixCmd) Build() (*Hotfix, error) {
	hotfix, err := store.SelectFirstHotfix(cmd.HotfixOptions)
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
	hotfixes, err := store.SelectHotfixes(cmd.HotfixOptions)
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
	// save hotfix
	entity := hotfix.Hotfix
	artifact := hotfix.HotfixArtifact
	entity.ArtifactArchs = strings.Join(artifact.ArtifactArchs, ",")
	entity.ArtifactEditions = strings.Join(artifact.ArtifactEditions, ",")
	entity.ArtifactTypes = strings.Join(artifact.ArtifactTypes, ",")
	if err := store.CreateOrUpdateHotfix(&entity); err != nil {
		return err
	}

	// save hotfixReleaseInfo
	releases := hotfix.ReleaseInfos
	for _, release := range releases {
		err := HotfixReleaseCmd{}.Save(release)
		if err != nil {
			return err
		}
	}

	return nil
}
