package model

import "tirelease/internal/entity"

type HotfixReleaseInfo struct {
	entity.HotfixReleaseInfo
	Issues    []Issue       `json:"issues,omitempty"`
	MasterPrs []PullRequest `json:"master_prs,omitempty"`
	BranchPrs []PullRequest `json:"branch_prs,omitempty"`
}
