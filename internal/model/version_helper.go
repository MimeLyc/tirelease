package model

import (
	"fmt"
	"strconv"
	"strings"
	"tirelease/commons/git"
	"tirelease/internal/entity"
)

// ====================================================
// ==================================================== Compose Function
func (version ReleaseVersion) ComposeVersionName() string {
	entityVersion := version
	if entityVersion.Addition == "" {
		return fmt.Sprintf("%d.%d.%d", entityVersion.Major, entityVersion.Minor, entityVersion.Patch)
	} else {
		return fmt.Sprintf("%d.%d.%d-%s", entityVersion.Major, entityVersion.Minor, entityVersion.Patch, entityVersion.Addition)
	}
}

func (version ReleaseVersion) ComposeVersionMinorName() string {
	return ComposeVersionMinorNameByNumber(version.Major, version.Minor)
}

func ComposeVersionMinorNameByNumber(major, minor int) string {
	return fmt.Sprintf("%d.%d", major, minor)
}

func ExtractVersionMinorName(versionName string) string {
	major, minor, _, _ := ComposeVersionAtom(versionName)
	return fmt.Sprintf("%d.%d", major, minor)
}

func (version ReleaseVersion) ComposeVersionBranch() string {
	return fmt.Sprintf("%s%d.%d", git.ReleaseBranchPrefix, version.Major, version.Minor)
}

func (version ReleaseVersion) ComposeVersionType() entity.ReleaseVersionType {
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
