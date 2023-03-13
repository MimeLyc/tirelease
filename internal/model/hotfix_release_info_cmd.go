package model

import (
	"strings"
	"tirelease/internal/entity"
	"tirelease/internal/store"
)

type HotfixReleaseCmd struct {
	HotfixReleaseInfoOptions *entity.HotfixReleaseInfoOptions
}

func (cmd HotfixReleaseCmd) BuildArray() ([]HotfixReleaseInfo, error) {
	entities, err := store.SelectHotfixReleaseInfos(cmd.HotfixReleaseInfoOptions)
	if err != nil {
		return nil, err
	}
	releases := entity.HotfixReleaseEntities(*entities)

	issueIds := releases.ExtractIssueIds()
	issues, err := IssueCmd{
		IssueOption: &entity.IssueOption{IssueIDs: issueIds},
	}.BuildArray()
	if err != nil {
		return nil, err
	}

	prIds := releases.ExtractAllPrIds()
	prs, err := PullRequestCmd{
		PROptions: &entity.PullRequestOption{PullRequestIDs: prIds},
	}.Build()
	if err != nil {
		return nil, err
	}

	result := make([]HotfixReleaseInfo, 0)
	for _, entity := range releases {
		release := HotfixReleaseInfo{
			HotfixReleaseInfo: entity,
		}

		issueIds = entity.ExtractIssueIds()
		release.Issues = filterIssuesByIssueIds(issues, issueIds)

		prIds = entity.ExtractMasterPrIds()
		release.MasterPrs = filterPrsByPrIds(prs, prIds)

		prIds = entity.ExtractBranchPrIds()
		release.BranchPrs = filterPrsByPrIds(prs, prIds)

		if entity.AssigneeEmail != "" {
			assignee, err := UserCmd{}.BuildByEmail(entity.AssigneeEmail)
			if err != nil {
				return nil, err
			}
			release.Assignee = assignee
		}

		result = append(result, release)
	}

	return result, nil
}

func (cmd HotfixReleaseCmd) Save(release HotfixReleaseInfo) error {
	entity := release.HotfixReleaseInfo
	entity.IssueIDs = strings.Join(release.ExtractIssueIds(), ",")
	entity.MasterPrIDs = strings.Join(release.ExtractMasterPrIds(), ",")
	entity.BranchPrIDs = strings.Join(release.ExtractBranchPrIds(), ",")

	return store.CreateOrUpdateHotfixReleaseInfo(&entity)
}
