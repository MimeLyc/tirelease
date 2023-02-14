package model

import (
	"tirelease/internal/entity"
)

type IssueVersionTriage struct {
	ID          int64
	Version     *ReleaseVersion
	Affect      entity.AffectResultResult
	Issue       *entity.Issue
	RelatedPrs  []PullRequest
	PickTriage  *pickTriageStateContext
	BlockTriage *BlockTriageStateContext

	Entity *entity.VersionTriage
	// All triage history of target issue for calculating block status and frontend render.
	HistoricalTriages *[]entity.VersionTriage
}

func (versionTriage IssueVersionTriage) MapToEntity() entity.VersionTriage {
	result := entity.VersionTriage{}
	if versionTriage.Entity != nil {
		result = *versionTriage.Entity
	}

	result.ID = versionTriage.ID
	result.VersionName = versionTriage.Version.Name
	result.IssueID = versionTriage.Issue.IssueID
	result.TriageResult = ParseToEntityPickTriage(versionTriage.PickTriage.State.StateText)
	result.BlockVersionRelease = ParseToEntityBlockTriage(versionTriage.BlockTriage.State.StateText)
	return result

}

func (versionTriage IssueVersionTriage) GetMergeStatus() entity.VersionTriageMergeStatus {
	if len(versionTriage.RelatedPrs) == 0 {
		return entity.VersionTriageMergeStatusPr
	}

	allMerge := true
	closeNums := 0
	for _, pr := range versionTriage.RelatedPrs {
		// PR state is closed when it's closed/cancelled or merged.
		// PR is closed/cancelled when PR state is "closed" and pr is not merged
		if pr.IsClosed() {
			closeNums++
			continue
		}

		//TODO: 当前存在approve成功hook到git，但是数据库中状态不一致的问题
		// 这里先兼容该情况，认为merge后的pr都是已approve过的，待重新设计状态机后修改逻辑
		if pr.IsMerged() {
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

	if closeNums == len(versionTriage.RelatedPrs) {
		return entity.VersionTriageMergeStatusPr
	}
	if allMerge {
		return entity.VersionTriageMergeStatusMerged
	} else {
		return entity.VersionTriageMergeStatusCITesting
	}
}

func (versionTriage *IssueVersionTriage) TriagePickStatus(status entity.VersionTriageResult) error {
	toStateText := ParseFromEntityPickTriage(status)
	pickTriage := versionTriage.PickTriage

	_, err := pickTriage.Trans(toStateText)

	return err
}

func (versionTriage *IssueVersionTriage) TriageBlockStatus(status entity.BlockVersionReleaseResult) error {
	toStateText := ParseFromEntityBlockTriage(status)
	blockTriage := versionTriage.BlockTriage

	_, err := blockTriage.Trans(toStateText)
	return err
}

func (versionTriage *IssueVersionTriage) ForceTriagePickStatus(status entity.VersionTriageResult) error {
	toStateText := ParseFromEntityPickTriage(status)
	pickTriage := versionTriage.PickTriage
	pickTriage.IsForce = true

	_, err := pickTriage.Trans(toStateText)

	return err
}
