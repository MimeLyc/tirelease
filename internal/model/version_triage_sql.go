package model

import (
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

func CreateOrUpdateVersionTriageInfo(triage *IssueVersionTriage, updatedVars ...entity.VersionTriageUpdatedVar) error {
	versionTriageDO := triage.MapToEntity()
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
//  such as the version name of version triage.
func SelectActiveIssueVersionTriage(versionName, issueID string) (*IssueVersionTriage, error) {
	// find version
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

	relatedPrs, err := GetRelatedPrs(releaseBranch, issueID)

	if err != nil {
		return nil, err
	}

	storedVersionTriage, err := selectMinorVersionTriage(releaseVersion.Name, issueID)
	if err != nil {
		return nil, err
	}

	affect := SelectVersionAffectResult(issueID, releaseVersion.ComposeVersionMinorName())

	blockTriage, err := NewBlockTriageStateContext(
		ParseFromEntityBlockTriage(storedVersionTriage.BlockVersionRelease),
		issue,
		releaseVersion,
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
		ID:          storedVersionTriage.ID,
		Version:     releaseVersion,
		Affect:      affect,
		Issue:       issue,
		RelatedPrs:  relatedPrs,
		PickTriage:  pickTriage,
		BlockTriage: blockTriage,
	}

	return &versionTriage, nil
}
