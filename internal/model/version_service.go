package model

import (
	"fmt"
	"strconv"
	"strings"
	"tirelease/commons/git"
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

// Select issue triages can be triaged under the version.
// If the version is finished, the result is equal to func: SelectHistoryIssueTriages
//   Else the triages will includes:
//       1. Issues affect the version but not triaged
//       2. Issues in the follow status: unknown, later, accepted, accepted(frozen)
func (version *ReleaseVersion) SelectCandidateIssueTriages() ([]IssueVersionTriage, error) {
	if !version.IsActive() {
		return version.SelectHistoryIssueTriages()
	}

	// Select active triages under version
	minorVersion := version.ComposeVersionMinorName()
	affectOption := &entity.IssueAffectOption{
		AffectVersion: minorVersion,
		AffectResult:  entity.AffectResultResultYes,
	}
	issueAffects, err := repository.SelectIssueAffect(affectOption)
	if err != nil {
		return nil, err
	}

	triages, err := selectMinorVersionTriages(version.Major, version.Minor)
	if err != nil {
		return nil, err
	}

	// Compose candidate issue triages: affects - history triages
	// TODO get the default block value of not triaged issues.
	triages = getCandidateTriages(version.Name, issueAffects, &triages)

	return composeVersionTriages(&triages, version, issueAffects)
}

// Select issue triages result under the version.
// * The issues affect the version but have not been triaged will not be selected.
func (version *ReleaseVersion) SelectHistoryIssueTriages() ([]IssueVersionTriage, error) {
	minorVersion := version.ComposeVersionMinorName()

	triages, err := repository.SelectVersionTriage(
		&entity.VersionTriageOption{
			VersionName: version.Name,
		},
	)
	if err != nil {
		return nil, err
	}

	// TODO: add a mechanism to deal with the triaged but no longer affected issues
	affects, err := repository.SelectIssueAffect(
		&entity.IssueAffectOption{
			AffectVersion: minorVersion,
		},
	)
	if err != nil {
		return nil, err
	}

	return composeVersionTriages(triages, version, affects)
}

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

func (version ReleaseVersion) ComposeVersionBranch() string {
	return fmt.Sprintf("%s%d.%d", git.ReleaseBranchPrefix, version.Major, version.Minor)
}

func ComposeVersionMinorNameByNumber(major, minor int) string {
	return fmt.Sprintf("%d.%d", major, minor)
}

func ExtractVersionMinorName(versionName string) string {
	major, minor, _, _ := ComposeVersionAtom(versionName)
	return fmt.Sprintf("%d.%d", major, minor)
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

func InitVersionStatus(vt entity.ReleaseVersionType) entity.ReleaseVersionStatus {
	if entity.ReleaseVersionTypePatch == vt {
		return entity.ReleaseVersionStatusPlanned
	} else {
		return entity.ReleaseVersionStatusUpcoming
	}
}

func getCandidateTriages(version string, issueAffects *[]entity.IssueAffect, issueTriages *[]entity.VersionTriage) []entity.VersionTriage {
	versionTriages := make([]entity.VersionTriage, 0)
	for i := range *issueAffects {
		issueAffect := (*issueAffects)[i]
		find := false

		// see if there is triage history
		for j := range *issueTriages {
			versionTriage := (*issueTriages)[j]
			if issueAffect.IssueID != versionTriage.IssueID {
				continue
			}
			find = true

			// if the issue is triage on lower version, remove it
			if versionTriage.TriageResult == entity.VersionTriageResultReleased ||
				versionTriage.TriageResult == entity.VersionTriageResultWontFix ||
				versionTriage.TriageResult == entity.VersionTriageResultLater {
				if version != versionTriage.VersionName {
					continue
				}
			}

			versionTriages = append(versionTriages, versionTriage)
		}

		if !find {
			versionTriage := entity.VersionTriage{
				IssueID:      issueAffect.IssueID,
				VersionName:  version,
				TriageResult: entity.VersionTriageResultUnKnown,
			}
			versionTriages = append(versionTriages, versionTriage)
		}
	}
	return versionTriages
}
