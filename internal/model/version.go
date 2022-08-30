package model

import (
	"errors"
	"fmt"
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

type ReleaseVersion struct {
	*entity.ReleaseVersion
}

func (version *ReleaseVersion) IsFrozen() bool {
	return version.Status == entity.ReleaseVersionStatusFrozen
}

func ComposeActiveReleaseVersion(name string) (*ReleaseVersion, error) {
	// release_version option
	shortType := ComposeVersionShortType(name)
	major, minor, patch, _ := ComposeVersionAtom(name)
	option := &entity.ReleaseVersionOption{}
	if shortType == entity.ReleaseVersionShortTypeMinor {
		option.Major = major
		option.Minor = minor
		option.StatusList = []entity.ReleaseVersionStatus{entity.ReleaseVersionStatusUpcoming, entity.ReleaseVersionStatusFrozen}
		option.ShortType = entity.ReleaseVersionShortTypeMinor
	} else if shortType == entity.ReleaseVersionShortTypePatch || shortType == entity.ReleaseVersionShortTypeHotfix {
		option.Major = major
		option.Minor = minor
		option.Patch = patch
		option.ShortType = entity.ReleaseVersionShortTypePatch
	} else {
		return nil, errors.New(fmt.Sprintf("SelectReleaseVersionActive params invalid: %+v failed", name))
	}

	// find version
	entityVersion, err := repository.SelectReleaseVersionLatest(option)
	if err != nil {
		return nil, err
	}
	releaseVersion := ReleaseVersion{
		entityVersion,
	}
	return &releaseVersion, nil

}
