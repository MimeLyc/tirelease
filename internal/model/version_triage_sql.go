package model

import (
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

func CreateOrUpdateVersionTriageInfo(triage *IssueVersionTriage, updatedVars ...entity.VersionTriageUpdatedVar) error {
	versionTriageDO := triage.MapToEntity()
	if len(updatedVars) == 0 {
		updatedVars = []entity.VersionTriageUpdatedVar{
			entity.VersionTriageUpdatedVarTriageResult,
			entity.VersionTriageUpdatedVarBlockRelease,
		}
	}
	return repository.CreateOrUpdateVersionTriage(&versionTriageDO, updatedVars...)
}

func SelectVersionAffectResult(issueID, minorVersionName string) entity.AffectResultResult {
	affect, err := repository.SelectIssueAffectUnique(&entity.IssueAffectOption{
		AffectVersion: minorVersionName,
		IssueID:       issueID,
	})

	if err != nil || affect == nil {
		return entity.AffectResultResultUnKnown
	}

	return affect.AffectResult
}

// Select the issue triaged under target version.
// There should only be one issue triage result under a minor version.
func selectMinorVersionTriage(versionName, issueID string) (*entity.VersionTriage, error) {
	storedVersionTriages, err := repository.SelectVersionTriage(
		&entity.VersionTriageOption{
			IssueID: issueID,
		})
	if err != nil {
		return nil, err
	}

	minorVersionName := ExtractVersionMinorName(versionName)

	for _, versionTriage := range *storedVersionTriages {
		if ExtractVersionMinorName(versionTriage.VersionName) == minorVersionName {
			versionTriage.VersionName = versionName
			return &versionTriage, nil
		}
	}

	return &entity.VersionTriage{
		VersionName: versionName,
		IssueID:     issueID,
	}, nil

}

// Compose single triage info towards single issue.
// **Will aotumatically set the related status to the active release version.**
//
//	such as the version name of version triage.
func SelectActiveIssueVersionTriage(versionName, issueID string) (*IssueVersionTriage, error) {
	// Find active patch version under target minor version.
	// Will replace the original version of triage under the minor version.
	releaseVersion, err := SelectActiveReleaseVersion(versionName)
	if err != nil {
		return nil, err
	}
	releaseBranch := releaseVersion.ReleaseBranch

	issue, err := repository.SelectIssueUnique(&entity.IssueOption{
		IssueID: issueID,
	})
	if err != nil {
		return nil, err
	}

	relatedPrs, err := SelectRelatedPrs(releaseBranch, issueID)

	if err != nil {
		return nil, err
	}

	storedVersionTriage, err := selectMinorVersionTriage(releaseVersion.Name, issueID)
	if err != nil {
		return nil, err
	}

	affect := SelectVersionAffectResult(issueID, releaseVersion.ComposeVersionMinorName())

	historyTriages, err := repository.SelectVersionTriage(
		&entity.VersionTriageOption{
			IssueID: issueID,
		},
	)
	if err != nil {
		return nil, err
	}

	blockTriage, err := NewBlockTriageStateContext(
		ParseFromEntityBlockTriage(storedVersionTriage.BlockVersionRelease),
		issue,
		releaseVersion,
		historyTriages,
	)

	if err != nil {
		return nil, err
	}

	pickTriage, err := NewPickTriageStateContext(
		ParseFromEntityPickTriage(storedVersionTriage.TriageResult),
		issue, releaseVersion, relatedPrs,
	)
	if err != nil {
		return nil, err
	}

	versionTriage := IssueVersionTriage{
		ID:                storedVersionTriage.ID,
		Version:           releaseVersion,
		Affect:            affect,
		Issue:             issue,
		RelatedPrs:        relatedPrs,
		PickTriage:        pickTriage,
		BlockTriage:       blockTriage,
		Entity:            storedVersionTriage,
		HistoricalTriages: historyTriages,
	}

	return &versionTriage, nil
}

func selectMinorVersionTriages(major, minor int) ([]entity.VersionTriage, error) {
	// select all triaged list under this minor version
	versionOption := &entity.ReleaseVersionOption{
		Major:     major,
		Minor:     minor,
		ShortType: entity.ReleaseVersionShortTypeMinor,
	}
	releaseVersions, err := repository.SelectReleaseVersion(versionOption)
	if err != nil {
		return nil, err
	}
	versions := make([]string, 0)
	for i := range *releaseVersions {
		versions = append(versions, (*releaseVersions)[i].Name)
	}

	versionTriageOption := &entity.VersionTriageOption{
		VersionNameList: versions,
	}
	versionTriageData, err := repository.SelectVersionTriage(versionTriageOption)
	if err != nil {
		return nil, err
	}
	return *versionTriageData, nil
}

func composeVersionTriages(triages *[]entity.VersionTriage, version *ReleaseVersion,
	affects *[]entity.IssueAffect) ([]IssueVersionTriage, error) {

	issueIDs := extractIssueIDsFromTriage(*triages)
	issues, err := SelectIssues(
		&entity.IssueOption{
			IssueIDs: issueIDs,
		},
	)
	if err != nil {
		return nil, err
	}

	history, err := repository.SelectVersionTriage(
		&entity.VersionTriageOption{
			IssueIDs: issueIDs,
		},
	)
	if err != nil {
		return nil, err
	}

	relations, err := SelectIssuePrRelationsByVersion(version.Major, version.Minor,
		entity.IssueOption{
			IssueIDs: issueIDs,
		}, false,
	)
	if err != nil {
		return nil, err
	}

	result := make([]IssueVersionTriage, 0)

	for _, triage := range *triages {
		triage := triage

		affect := FilterAffectByIssueIDandMinorVersion(*affects, triage.IssueID, version.ComposeVersionMinorName())
		issue := FilterIssueByID(*issues, triage.IssueID)

		relation := FilterIssuePrRelationByIssueAndVersion(relations, issue.IssueID, version.Major, version.Minor)
		relatedPrs := make([]entity.PullRequest, 0)
		if relation != nil {
			relatedPrs = relation.RelatedPrs
		}

		historyTriages := filterTriagesByIssue(*history, triage.IssueID)

		pickTriage, _ := NewPickTriageStateContext(StateText(triage.TriageResult), issue, version, relatedPrs)
		blockTriage, _ := NewBlockTriageStateContext(StateText(triage.BlockVersionRelease), issue, version, &historyTriages)

		affectResult := entity.AffectResultResultUnKnown
		if affect != nil {
			affectResult = affect.AffectResult
		}
		issueVersionTriage := IssueVersionTriage{
			ID:                triage.ID,
			Version:           version,
			Affect:            affectResult,
			Issue:             issue,
			RelatedPrs:        relatedPrs,
			PickTriage:        pickTriage,
			BlockTriage:       blockTriage,
			Entity:            &triage,
			HistoricalTriages: &historyTriages,
		}
		result = append(result, issueVersionTriage)
	}

	return result, nil
}

func filterTriagesByIssue(triages []entity.VersionTriage, issueID string) []entity.VersionTriage {
	result := make([]entity.VersionTriage, 0)
	for _, triage := range triages {
		if triage.IssueID == issueID {
			result = append(result, triage)
		}
	}
	return result
}

func SelectAllTriagesByIssue(issue entity.Issue) ([]IssueVersionTriage, error) {
	affects, err := repository.SelectIssueAffect(
		&entity.IssueAffectOption{
			IssueID:      issue.IssueID,
			AffectResult: entity.AffectResultResultYes,
		},
	)
	if err != nil {
		return nil, err
	}

	triages, err := repository.SelectVersionTriage(
		&entity.VersionTriageOption{
			IssueID: issue.IssueID,
		},
	)
	if err != nil {
		return nil, err
	}

	relations, err := repository.SelectIssuePrRelation(
		&entity.IssuePrRelationOption{
			IssueID: issue.IssueID,
		},
	)
	if err != nil {
		return nil, err
	}

	prs, err := selectRelatedPullRequests(*relations)
	if err != nil {
		return nil, err
	}

	result := make([]IssueVersionTriage, 0)
	for _, affect := range *affects {
		affect := affect
		minorName := affect.AffectVersion

		isTriaged := false
		triage := entity.VersionTriage{
			IssueID: issue.IssueID,
		}
		for j := range *triages {
			versionTriage := (*triages)[j]
			major, minor, _, _ := ComposeVersionAtom(versionTriage.VersionName)
			if ComposeVersionMinorNameByNumber(major, minor) == minorName {
				triage = versionTriage
				isTriaged = true
			}
		}

		// If the issue has been triaged, get the target version
		// Else get the active one.
		version := &ReleaseVersion{}
		if isTriaged {
			version, err = SelectReleaseVersion(triage.VersionName)
		} else {
			version, err = SelectActiveReleaseVersion(minorName)
		}
		if _, ok := err.(repository.VersionNotFoundError); ok {
			continue
		}
		if err != nil {
			return nil, err
		}
		triage.VersionName = version.Name

		ComposeVersionMinorNameByNumber(version.Major, version.Minor)
		branch := version.ComposeVersionBranch()

		relatedPrs := make([]entity.PullRequest, 0)
		for _, pr := range prs {
			pr := pr
			if pr.BaseBranch == branch {
				relatedPrs = append(relatedPrs, pr)
			}
		}

		pickTriage, _ := NewPickTriageStateContext(StateText(triage.TriageResult), &issue, version, relatedPrs)
		blockTriage, _ := NewBlockTriageStateContext(StateText(triage.BlockVersionRelease), &issue, version, triages)

		issueVersionTriage := IssueVersionTriage{
			ID:                triage.ID,
			Version:           version,
			Affect:            affect.AffectResult,
			Issue:             &issue,
			RelatedPrs:        relatedPrs,
			PickTriage:        pickTriage,
			BlockTriage:       blockTriage,
			Entity:            &triage,
			HistoricalTriages: triages,
		}
		result = append(result, issueVersionTriage)
	}

	return result, nil
}
