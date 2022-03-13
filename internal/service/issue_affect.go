package service

import (
	"fmt"
	"strings"

	"tirelease/commons/git"
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

// ============================================================================
// ============================================================================ Restful API(From UI) Handler
func CreateOrUpdateIssueAffect(issueAffect *entity.IssueAffect) error {
	// create or update
	err := repository.CreateOrUpdateIssueAffect(issueAffect)
	if err != nil {
		return err
	}

	// result operation
	err = OperateIssueAffectResult(issueAffect)
	if err != nil {
		return err
	}

	return nil
}

func OperateIssueAffectResult(issueAffect *entity.IssueAffect) error {
	// param protection
	if issueAffect.AffectResult == "" {
		return nil
	}

	// operate git label
	affectLabel := fmt.Sprintf(git.AffectsLabel, issueAffect.AffectVersion)
	mayAffectLabel := fmt.Sprintf(git.MayAffectsLabel, issueAffect.AffectVersion)
	if issueAffect.AffectResult != entity.AffectResultResultUnKnown {
		err := RemoveLabelByIssueID(issueAffect.IssueID, mayAffectLabel)
		if err != nil {
			return err
		}
	}
	if issueAffect.AffectResult == entity.AffectResultResultYes {
		err := AddLabelByIssueID(issueAffect.IssueID, affectLabel)
		if err != nil {
			return err
		}
	}
	if issueAffect.AffectResult == entity.AffectResultResultNo {
		err := RemoveLabelByIssueID(issueAffect.IssueID, affectLabel)
		if err != nil {
			return err
		}
	}

	// operate cherry-pick record: select latest version & insert cherry-pick
	if issueAffect.AffectResult == entity.AffectResultResultYes {
		releaseVersionOption := &entity.ReleaseVersionOption{
			FatherReleaseVersionName: issueAffect.AffectVersion,
			Status:                   entity.ReleaseVersionStatusOpen,
			// Type:                     entity.ReleaseVersionTypePatch,
		}
		releaseVersion, err := repository.SelectReleaseVersionLatest(releaseVersionOption)
		if err != nil {
			return err
		}
		versionTriage := &entity.VersionTriage{
			IssueID:      issueAffect.IssueID,
			VersionName:  releaseVersion.Name,
			TriageResult: entity.VersionTriageResultUnKnown,
		}
		_, err = CreateOrUpdateVersionTriageInfo(versionTriage)
		if err != nil {
			return err
		}
	}

	return nil
}

// ============================================================================
// ============================================================================ Compose From DataBase & Return To UI
func ComposeIssueAffectWithIssueID(issueID string) (*[]entity.IssueAffect, error) {
	// Select Exist Issue Affect
	issueAffectOption := &entity.IssueAffectOption{
		IssueID: issueID,
	}
	issueAffects, err := repository.SelectIssueAffect(issueAffectOption)
	if err != nil {
		return nil, err
	}

	// Implement New Issue Affect
	releaseVersionOption := &entity.ReleaseVersionOption{
		Type:   entity.ReleaseVersionTypeMinor,
		Status: entity.ReleaseVersionStatusOpen,
	}
	releaseVersions, err := repository.SelectReleaseVersion(releaseVersionOption)
	if nil != err {
		return nil, err
	}
	for _, releaseVersion := range *releaseVersions {
		var isExist bool = false
		for _, issueAffect := range *issueAffects {
			if issueAffect.AffectVersion == releaseVersion.Name {
				isExist = true
				break
			}
		}

		if !isExist {
			newAffect := &entity.IssueAffect{
				IssueID:       issueID,
				AffectVersion: releaseVersion.Name,
				AffectResult:  entity.AffectResultResultUnKnown,
			}
			(*issueAffects) = append((*issueAffects), *newAffect)
		}
	}
	return issueAffects, nil
}

// ============================================================================
// ============================================================================ Compose From Webhook (Save To DataBase)
func ComposeIssueAffectWithIssueV4(issue *git.IssueField) (*[]entity.IssueAffect, error) {
	if nil == issue || len(issue.Labels.Nodes) == 0 {
		return nil, nil
	}

	issueAffects := make([]entity.IssueAffect, 0)
	for _, label := range issue.Labels.Nodes {
		labelName := string(label.Name)
		if strings.HasPrefix(labelName, git.AffectsPrefixLabel) {
			version := strings.Replace(labelName, git.AffectsPrefixLabel, "", -1)
			issueAffect := entity.IssueAffect{
				IssueID:       issue.ID.(string),
				AffectVersion: version,
				AffectResult:  entity.AffectResultResultYes,
			}
			issueAffects = append(issueAffects, issueAffect)
		}
	}
	return &issueAffects, nil
}
