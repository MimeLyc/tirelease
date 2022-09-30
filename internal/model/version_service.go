package model

import (
	"fmt"
	"strconv"
	"strings"
	"tirelease/commons/git"
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

// Select issue triages result under the version.
// * The issues affect the version but have not been triaged will also be selected.
func (version *ReleaseVersion) SelectIssueTriages() ([]IssueVersionTriage, error) {
	minorVersion := version.ComposeVersionMinorName()

	triages, err := repository.SelectVersionTriage(
		&entity.VersionTriageOption{
			VersionName: version.Name,
		},
	)
	if err != nil {
		return nil, err
	}

	issueIDs := extractIssueIDsFromTriage(*triages)
	issues, err := repository.SelectIssue(
		&entity.IssueOption{
			IssueIDs: issueIDs,
		},
	)
	if err != nil {
		return nil, err
	}
	issuePrRelations, err := SelectIssuePrRelations(version.Major, version.Minor,
		entity.IssueOption{
			IssueIDs: issueIDs,
		}, false,
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

	result := make([]IssueVersionTriage, 0)
	for _, triage := range *triages {
		triage := triage
		affect := FilterAffectByIssueIDandMinorVersion(*affects, triage.IssueID, minorVersion)
		issue := FilterIssueByID(*issues, triage.IssueID)

		relation := FilterIssuePrRelationByIssueAndVersion(issuePrRelations, issue.IssueID, version.Major, version.Minor)
		relatedPrs := make([]entity.PullRequest, 0)
		if relation != nil {
			relatedPrs = relation.RelatedPrs
		}

		pickTriage, _ := NewPickTriageStateContext(StateText(triage.TriageResult), issue, version, relatedPrs)
		blockTriage, _ := NewBlockTriageStateContext(StateText(triage.BlockVersionRelease), issue, version)
		affectResult := entity.AffectResultResultUnKnown
		if affect != nil {
			affectResult = affect.AffectResult
		}
		issueVersionTriage := IssueVersionTriage{
			ID:          triage.ID,
			Version:     version,
			Affect:      affectResult,
			Issue:       issue,
			RelatedPrs:  relatedPrs,
			PickTriage:  pickTriage,
			BlockTriage: blockTriage,
		}
		result = append(result, issueVersionTriage)
	}

	return result, nil
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
