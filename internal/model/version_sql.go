package model

import (
	"errors"
	"fmt"
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

func CreateNextVersionIfNotExist(preVersion *ReleaseVersion) (*ReleaseVersion, error) {
	major, minor, patch, _ := ComposeVersionAtom(preVersion.Name)

	option := &entity.ReleaseVersionOption{
		Major:     major,
		Minor:     minor,
		Patch:     patch + 1,
		ShortType: entity.ReleaseVersionShortTypePatch,
	}
	version := ReleaseVersion{}

	versionEntity, err := repository.SelectReleaseVersionLatest(option)

	if nil == err && nil != versionEntity {
		version = Parse2ReleaseVersion(*versionEntity)
	} else {
		version = Parse2ReleaseVersion(entity.ReleaseVersion{
			Major: major,
			Minor: minor,
			Patch: patch + 1,
		})

		// TODO
		version, err = CreateReleaseVersion(version)
		if nil != err {
			return nil, err
		}
	}
	return &version, nil
}

func CreateReleaseVersion(releaseVersion ReleaseVersion) (ReleaseVersion, error) {
	releaseVersion, _ = initReleaseVersion(releaseVersion)
	err := repository.CreateReleaseVersion(releaseVersion.ReleaseVersion)
	if nil != err {
		return releaseVersion, err
	}

	return releaseVersion, nil
}

// Select last patch version of param version.
// While the patch number is 0, return nil
func SelectLastPatchVersion(version ReleaseVersion) (*ReleaseVersion, error) {
	if version.Patch == 0 {
		return nil, nil
	}

	patch := version.Patch - 1

	versionEntity, err := repository.SelectReleaseVersionLatest(
		&entity.ReleaseVersionOption{
			Major: version.Major,
			Minor: version.Minor,
			Patch: patch,
		},
	)
	if err != nil {
		return nil, err
	}
	releaseVersion := Parse2ReleaseVersion(*versionEntity)
	return &releaseVersion, nil
}

func SelectActiveReleaseVersion(name string) (*ReleaseVersion, error) {
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
	releaseVersion := Parse2ReleaseVersion(*entityVersion)
	return &releaseVersion, nil

}
