package entity

import "time"

type HotfixReleaseInfo struct {
	ID int64 `json:"id"`

	CreateTime time.Time `json:"create_time,omitempty"`
	UpdateTime time.Time `json:"update_time,omitempty"`

	HotfixName string `json:"hotfix_name,omitempty"`
	// Virtual foreign key for mapping entity.Repo
	RepoFullName string `json:"repo_full_name,omitempty"`

	BasedReleaseVersion string `json:"based_release_version,omitempty"`
	BasedCommitSha      string `json:"based_commit_sha,omitempty"`

	// Virtual foreign keys for mapping entity.Issue
	IssueIDs string `json:"issue_ids,omitempty"`

	// Virtual foreign keys for mapping entity.PullRequest
	MasterPrIDs string `json:"master_pr_ids,omitempty"`

	BuildID      string `json:"build_id,omitempty"`
	BuildFeature string `json:"build_feature,omitempty"`
}

// DB-Table
func (HotfixReleaseInfo) TableName() string {
	return "hotfix_release_info"
}
