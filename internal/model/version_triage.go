package model

import (
	"tirelease/internal/entity"
)

type IssueVersionTriage struct {
	ID          int64
	Version     *ReleaseVersion
	Affect      entity.AffectResultResult
	Issue       *entity.Issue
	RelatedPrs  []entity.PullRequest
	PickTriage  *PickTriageStateContext
	BlockTriage *BlockTriageStateContext
}

func (versionTriage IssueVersionTriage) MapToEntity() entity.VersionTriage {
	return entity.VersionTriage{
		ID:                  versionTriage.ID,
		VersionName:         versionTriage.Version.Name,
		IssueID:             versionTriage.Issue.IssueID,
		TriageResult:        ParseToEntityPickTriage(versionTriage.PickTriage.State.StateText),
		BlockVersionRelease: ParseToEntityBlockTriage(versionTriage.BlockTriage.State.StateText),
	}
}

func (versionTriage IssueVersionTriage) GetMergeStatus() entity.VersionTriageMergeStatus {
	if len(versionTriage.RelatedPrs) == 0 {
		return entity.VersionTriageMergeStatusPr
	}

	allMerge := true
	closeNums := 0
	for _, pr := range versionTriage.RelatedPrs {
		if pr.State == "closed" {
			closeNums++
			continue
		}

		//TODO: 当前存在approve成功hook到git，但是数据库中状态不一致的问题
		// 这里先兼容该情况，认为merge后的pr都是已approve过的，待重新设计状态机后修改逻辑
		if pr.Merged {
			continue
		}

		if !pr.CherryPickApproved {
			return entity.VersionTriageMergeStatusApprove
		} else if !pr.AlreadyReviewed {
			return entity.VersionTriageMergeStatusReview
		} else if !pr.Merged {
			allMerge = false
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
