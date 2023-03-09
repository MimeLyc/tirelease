package service

import (
	"strings"
	"tirelease/commons/git"
	"tirelease/commons/utils"
	"tirelease/internal/dto"
	"tirelease/internal/entity"
	"tirelease/internal/model"
	"tirelease/internal/store"

	"github.com/pkg/errors"
)

// Create Or Update version triage info, including block triage and pick triage.
// Attention: the version in **versionTriage** param will always be a minor version.
//
//	so the triage result will be moved to the latest patch version automatically.
func SaveVersionTriageInfo(versionTriage *entity.VersionTriage, updatedVars ...entity.VersionTriageUpdatedVar) (*dto.VersionTriageInfo, error) {
	issueVersionTriage, err := model.SelectActiveIssueVersionTriage(versionTriage.VersionName, versionTriage.IssueID)
	if err != nil {
		return nil, err
	}

	if utils.Contains(updatedVars, entity.VersionTriageUpdatedVarTriageResult) {
		issueVersionTriage.TriagePickStatus(versionTriage.TriageResult)
		versionTriage.TriageResult = model.ParseToEntityPickTriage(issueVersionTriage.PickTriage.State.StateText)
	}

	versionTriage.VersionName = issueVersionTriage.Version.Name
	// set default Block Triage status
	versionTriage.BlockVersionRelease = model.ParseToEntityBlockTriage(issueVersionTriage.BlockTriage.State.StateText)

	model.CreateOrUpdateVersionTriageInfo(issueVersionTriage, updatedVars...)

	// return
	return &dto.VersionTriageInfo{
		ReleaseVersion: issueVersionTriage.Version.ReleaseVersion,
		IsFrozen:       issueVersionTriage.Version.IsFrozen(),
		IsAccept:       issueVersionTriage.PickTriage.IsAccept(),

		VersionTriage:            versionTriage,
		VersionTriageMergeStatus: issueVersionTriage.GetMergeStatus(),
		// deprecated: IssueRelationInfo in the related API is not used.
		IssueRelationInfo: nil,
	}, nil
}

func ForceTriageIssue(version, owner, repo string, number int, triageResult entity.VersionTriageResult) error {
	issue, err := model.IssueCmd{}.BuildByNumber(owner, repo, number)
	if err != nil {
		return err
	}

	return issue.ForcePickTriage(version, triageResult)
}

// Create Or Update batch triage info
// Now the method **only includes** pick triage logic.
// Attention: the version in **versionTriage** param will always be a minor version.
//
//	so the triage result will be moved to the latest patch version automatically.
func CreateOrUpdateIssueTriages(triages *[]entity.VersionTriageOption) error {
	for _, modifiedTriage := range *triages {
		triage, err := model.SelectActiveIssueVersionTriage(modifiedTriage.VersionName, modifiedTriage.IssueID)
		if err != nil {
			return err
		}

		originalTriage := model.ParseToEntityPickTriage(triage.PickTriage.GetStateText())
		if originalTriage == modifiedTriage.TriageResult {
			continue
		}

		triage.TriagePickStatus(modifiedTriage.TriageResult)

		err = model.CreateOrUpdateVersionTriageInfo(triage, entity.VersionTriageUpdatedVarTriageResult, entity.VersionTriageUpdatedVarVersion)
		if err != nil {
			return err
		}
	}

	return nil
}

// Hook github api to change the `approve` related labels.
func ChangePrApprovedLabel(prId string, isFrozen, isAccept bool) error {
	if !isFrozen && isAccept {
		err := RemoveLabelByPullRequestID(prId, git.NotCheryyPickLabel)
		if err != nil {
			return err
		}

		err = AddLabelByPullRequestID(prId, git.CherryPickLabel)
		if err != nil {
			return err
		}
	} else {
		err := RemoveLabelByPullRequestID(prId, git.CherryPickLabel)
		if err != nil {
			return err
		}

		err = AddLabelByPullRequestID(prId, git.NotCheryyPickLabel)
		if err != nil {
			return err
		}
	}

	return nil
}

func FindVersionTriageInfo(query *dto.VersionTriageInfoQuery) (*dto.VersionTriageInfoWrap, *entity.ListResponse, error) {
	version, err := model.SelectReleaseVersion(query.Version)
	if err != nil {
		return nil, nil, err
	}

	issueTriages, err := version.SelectCandidateIssueTriages()
	if err != nil {
		return nil, nil, err
	}

	// Map to versionTriageInfos for frontend display.
	versionTriageInfos := make([]dto.VersionTriageInfo, 0)
	for _, triage := range issueTriages {
		triage := triage
		info := MapToVersionTriageInfo(triage)
		versionTriageInfos = append(versionTriageInfos, info)
	}

	wrap := &dto.VersionTriageInfoWrap{
		ReleaseVersion:     version.ReleaseVersion,
		VersionTriageInfos: &versionTriageInfos,
	}
	response := &entity.ListResponse{
		Page:    query.VersionTriageOption.Page,
		PerPage: query.VersionTriageOption.PerPage,
	}
	response.CalcTotalPage()
	return wrap, response, nil
}

func UpdateVersionTriage(versionTriage *entity.VersionTriage) error {
	if versionTriage.ID == 0 {
		versionTriage.TriageResult = entity.VersionTriageResultUnKnown
		err := store.CreateVersionTriage(versionTriage)
		if err != nil {
			return err
		}
	} else {
		err := store.UpdateVersionTriage(versionTriage)
		if err != nil {
			return err
		}
	}
	return nil
}

func ComposeVersionTriageMergeStatus(relatedPrs []entity.PullRequest) entity.VersionTriageMergeStatus {
	if len(relatedPrs) == 0 {
		return entity.VersionTriageMergeStatusPr
	}

	allMerge := true
	closeNums := 0
	for _, pr := range relatedPrs {
		// PR state is closed when it's closed/cancelled or merged.
		// PR is closed/cancelled when PR state is "closed" and pr is not merged
		if pr.State == "closed" && !pr.Merged {
			closeNums++
			continue
		}

		//TODO: 当前存在approve成功hook到git，但是数据库中状态不一致的问题
		// 这里先兼容该情况，认为merge后的pr都是已approve过的，待重新设计状态机后修改逻辑
		if pr.Merged {
			continue
		} else {
			allMerge = false
		}

		if !pr.CherryPickApproved {
			return entity.VersionTriageMergeStatusApprove
		} else if !pr.AlreadyReviewed {
			return entity.VersionTriageMergeStatusReview
		}
	}

	if closeNums == len(relatedPrs) {
		return entity.VersionTriageMergeStatusPr
	}
	if allMerge {
		return entity.VersionTriageMergeStatusMerged
	} else {
		return entity.VersionTriageMergeStatusCITesting
	}
}

func ComposeVersionTriageUpcomingList(version string) ([]entity.VersionTriage, error) {
	// select all issue which may affect this minor version
	major, minor, _, _ := ComposeVersionAtom(version)
	minorVersion := ComposeVersionMinorNameByNumber(major, minor)
	affectOption := &entity.IssueAffectOption{
		AffectVersion: minorVersion,
		AffectResult:  entity.AffectResultResultYes,
	}
	issueAffects, err := store.SelectIssueAffect(affectOption)
	if err != nil {
		return nil, err
	}

	// select all triaged list under this minor version
	versionOption := &entity.ReleaseVersionOption{
		Major:     major,
		Minor:     minor,
		ShortType: entity.ReleaseVersionShortTypeMinor,
	}
	releaseVersions, err := store.SelectReleaseVersion(versionOption)
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
	versionTriageData, err := store.SelectVersionTriage(versionTriageOption)
	if err != nil {
		return nil, err
	}

	// compose: version_triage = affected - triaged
	versionTriages := make([]entity.VersionTriage, 0)
	for i := range *issueAffects {
		issueAffect := (*issueAffects)[i]
		find := false
		for j := range *versionTriageData {
			versionTriage := (*versionTriageData)[j]
			if issueAffect.IssueID != versionTriage.IssueID {
				continue
			}
			find = true

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

			// TODO refactor bellow logic of default Block value
			issueOption := entity.IssueOption{
				IssueID: issueAffect.IssueID,
			}
			issue, _ := store.SelectIssueUnique(&issueOption)
			if issue != nil && issue.SeverityLabel == git.SeverityCriticalLabel {
				versionTriage.BlockVersionRelease = entity.BlockVersionReleaseResultBlock
			}

			versionTriages = append(versionTriages, versionTriage)
		}
	}
	return versionTriages, nil
}

// Export history data (Only database operation, no remote operation)
func ExportHistoryVersionTriageInfo(info *dto.IssueRelationInfo, releaseVersions *[]entity.ReleaseVersion) error {
	// param check
	if info == nil || releaseVersions == nil {
		return errors.New("ExportHistoryVersionTriageInfo params invalid")
	}
	if info.PullRequests == nil || len(*info.PullRequests) == 0 {
		return nil
	}

	// insert version triage
	for i := range *info.PullRequests {
		pr := (*info.PullRequests)[i]
		if !pr.Merged || !pr.CherryPickApproved ||
			!strings.HasPrefix(pr.BaseBranch, git.ReleaseBranchPrefix) {
			continue
		}
		releaseBranch := string(pr.BaseBranch)
		branchVersion := strings.Replace(pr.BaseBranch, git.ReleaseBranchPrefix, "", -1)
		major, minor, _, _ := ComposeVersionAtom(branchVersion)

		// search version in time section
		// release version is already sorted desc
		for i := len(*releaseVersions) - 1; i >= 0; i-- {
			releaseVersion := (*releaseVersions)[i]
			if releaseVersion.Status != entity.ReleaseVersionStatusReleased {
				continue
			}
			if releaseVersion.Major != major || releaseVersion.Minor != minor || releaseVersion.ReleaseBranch != releaseBranch {
				continue
			}
			if releaseVersion.ActualReleaseTime.After(*(pr.MergeTime)) {
				versionTriage := &entity.VersionTriage{
					IssueID:      info.Issue.IssueID,
					VersionName:  releaseVersion.Name,
					TriageResult: entity.VersionTriageResultReleased,
					CreateTime:   *(pr.MergeTime),
					UpdateTime:   *(pr.MergeTime),
				}
				if err := store.CreateOrUpdateVersionTriage(versionTriage); err != nil {
					return err
				}

				issueAffect := &entity.IssueAffect{
					IssueID:       info.Issue.IssueID,
					AffectVersion: branchVersion,
					AffectResult:  entity.AffectResultResultYes,
					CreateTime:    *(pr.MergeTime),
					UpdateTime:    *(pr.MergeTime),
				}
				if err := store.CreateOrUpdateIssueAffect(issueAffect); err != nil {
					return err
				}

				break
			}
		}
	}

	return nil
}

func fillVersionTriageDefaultValue(issueRelationInfo *dto.IssueRelationInfo, versionTriage *entity.VersionTriage) {
	if len(*issueRelationInfo.VersionTriages) == 0 {
		if issueRelationInfo.Issue.SeverityLabel == git.SeverityCriticalLabel {
			versionTriage.BlockVersionRelease = entity.BlockVersionReleaseResultBlock
		}
	}
}
