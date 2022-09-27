package service

import (
	"fmt"
	"strconv"
	"strings"

	"tirelease/commons/git"
	"tirelease/internal/entity"
	"tirelease/internal/model"
	"tirelease/internal/repository"

	mapset "github.com/deckarep/golang-set"
	"github.com/pkg/errors"
)

func CreateReleaseVersion(releaseVersion *entity.ReleaseVersion) error {
	version := model.Parse2ReleaseVersion(*releaseVersion)
	_, err := model.CreateReleaseVersion(version)

	return err
}

// params
//  releaseVersion: the post request body from api which contain the full data of release version.
func UpdateReleaseVersion(releaseVersion *entity.ReleaseVersion) error {
	if releaseVersion.Type == entity.ReleaseVersionTypeHotfix {
		return nil
	}

	version := model.Parse2ReleaseVersion(*releaseVersion)
	return version.ChangeStatus(releaseVersion.Status)
}

func SelectReleaseVersion(option *entity.ReleaseVersionOption) (*[]entity.ReleaseVersion, error) {
	return repository.SelectReleaseVersion(option)
}

func SelectReleaseVersionMaintained() (*[]string, error) {
	option := &entity.ReleaseVersionOption{
		StatusList: []entity.ReleaseVersionStatus{
			entity.ReleaseVersionStatusUpcoming,
			entity.ReleaseVersionStatusFrozen,
		},
	}
	versions, err := repository.SelectReleaseVersion(option)
	if nil != err {
		return nil, err
	}

	set := mapset.NewSet()
	for _, version := range *versions {
		set.Add(ComposeVersionMinorName(&version))
	}
	var res []string
	for _, v := range set.ToSlice() {
		res = append(res, v.(string))
	}
	return &res, nil
}

// deprecated
func SelectReleaseVersionActive(name string) (*entity.ReleaseVersion, error) {
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
	releaseVersion, err := repository.SelectReleaseVersionLatest(option)
	if err != nil {
		return nil, err
	}
	return releaseVersion, nil
}

// ====================================================
// ==================================================== Compose Function
func ComposeVersionName(version *entity.ReleaseVersion) string {
	if version.Addition == "" {
		return fmt.Sprintf("%d.%d.%d", version.Major, version.Minor, version.Patch)
	} else {
		return fmt.Sprintf("%d.%d.%d-%s", version.Major, version.Minor, version.Patch, version.Addition)
	}
}

func ComposeVersionMinorName(version *entity.ReleaseVersion) string {
	return ComposeVersionMinorNameByNumber(version.Major, version.Minor)
}

func ComposeVersionMinorNameByNumber(major, minor int) string {
	return fmt.Sprintf("%d.%d", major, minor)
}

func ComposeVersionBranch(version *entity.ReleaseVersion) string {
	return fmt.Sprintf("%s%d.%d", git.ReleaseBranchPrefix, version.Major, version.Minor)
}

func ComposeVersionType(version *entity.ReleaseVersion) entity.ReleaseVersionType {
	if version.Addition != "" {
		return entity.ReleaseVersionTypeHotfix
	} else {
		if version.Patch != 0 {
			return entity.ReleaseVersionTypePatch
		} else {
			if version.Minor != 0 {
				return entity.ReleaseVersionTypeMinor
			} else {
				return entity.ReleaseVersionTypeMajor
			}
		}
	}
}

func ComposeVersionShortType(version string) entity.ReleaseVersionShortType {
	// todo: regexp later
	slice := strings.Split(version, "-")
	if len(slice) >= 2 {
		return entity.ReleaseVersionShortTypeHotfix
	}

	slice = strings.Split(slice[0], ".")
	if len(slice) == 3 {
		return entity.ReleaseVersionShortTypePatch
	}
	if len(slice) == 2 {
		return entity.ReleaseVersionShortTypeMinor
	}
	if len(slice) == 1 {
		return entity.ReleaseVersionShortTypeMajor
	}
	return entity.ReleaseVersionShortTypeUnKnown
}

func ComposeVersionAtom(version string) (major, minor, patch int, addition string) {
	major = 0
	minor = 0
	patch = 0
	addition = ""

	slice := strings.Split(version, "-")
	if len(slice) >= 2 {
		for _, v := range slice[1:] {
			addition += v
			if v != slice[len(slice)-1] {
				addition += "-"
			}
		}
	}

	slice = strings.Split(slice[0], ".")
	if len(slice) >= 1 {
		major, _ = strconv.Atoi(slice[0])
	}
	if len(slice) >= 2 {
		minor, _ = strconv.Atoi(slice[1])
	}
	if len(slice) >= 3 {
		patch, _ = strconv.Atoi(slice[2])
	}

	return major, minor, patch, addition
}

func ComposeVersionStatus(vt entity.ReleaseVersionType) entity.ReleaseVersionStatus {
	if entity.ReleaseVersionTypePatch == vt {
		return entity.ReleaseVersionStatusPlanned
	} else {
		return entity.ReleaseVersionStatusUpcoming
	}
}
