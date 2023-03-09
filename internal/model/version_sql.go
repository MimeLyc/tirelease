package model

import (
	"errors"
	"fmt"
	"tirelease/internal/entity"
	"tirelease/internal/store"
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

	versionEntity, err := store.SelectReleaseVersionLatest(option)

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
	err := store.CreateReleaseVersion(releaseVersion.ReleaseVersion)
	if nil != err {
		return releaseVersion, err
	}

	return releaseVersion, nil
}

// Select last patch version of param version.
// While the patch number is 0, return nil
func SelectPrePatchVersion(version ReleaseVersion) (*ReleaseVersion, error) {
	if version.Patch == 0 {
		return nil, nil
	}

	patch := version.Patch - 1

	versionEntity, err := store.SelectReleaseVersionLatest(
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

// Select active release version model under target minor version
func SelectActiveReleaseVersion(minorVersion string) (*ReleaseVersion, error) {
	// release_version option
	shortType := ComposeVersionShortType(minorVersion)
	major, minor, patch, _ := ComposeVersionAtom(minorVersion)
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
		return nil, errors.New(fmt.Sprintf("SelectReleaseVersionActive params invalid: %+v failed", minorVersion))
	}

	// find version
	entityVersion, err := store.SelectReleaseVersionLatest(option)
	if err != nil {
		return nil, err
	}
	releaseVersion := Parse2ReleaseVersion(*entityVersion)
	return &releaseVersion, nil

}

func SelectReleaseVersion(name string) (*ReleaseVersion, error) {
	// find version
	entityVersion, err := store.SelectReleaseVersionLatest(
		&entity.ReleaseVersionOption{
			Name: name,
		},
	)
	if err != nil {
		return nil, err
	}
	releaseVersion := Parse2ReleaseVersion(*entityVersion)
	return &releaseVersion, nil
}
