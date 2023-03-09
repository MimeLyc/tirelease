package model

import "tirelease/internal/entity"

type HotfixReleaseInfo struct {
	entity.HotfixReleaseInfo
	Issues    []Issue       `json:"issues,omitempty"`
	MasterPrs []PullRequest `json:"master_prs,omitempty"`
	BranchPrs []PullRequest `json:"branch_prs,omitempty"`
}

func (release HotfixReleaseInfo) ExtractIssueIds() []string {
	result := make([]string, 0)
	for _, issue := range release.Issues {
		result = append(result, issue.IssueID)
	}
	return result
}

func (release HotfixReleaseInfo) ExtractMasterPrIds() []string {
	result := make([]string, 0)
	for _, pr := range release.MasterPrs {
		result = append(result, pr.PullRequestID)
	}
	return result
}

func (release HotfixReleaseInfo) ExtractBranchPrIds() []string {
	result := make([]string, 0)
	for _, pr := range release.BranchPrs {
		result = append(result, pr.PullRequestID)
	}
	return result
}
