package entity

import (
	"strings"
	"time"
)

type HotfixReleaseInfo struct {
	ID int64 `json:"id"`

	CreateTime time.Time `json:"create_time,omitempty"`
	UpdateTime time.Time `json:"update_time,omitempty"`

	HotfixName string `json:"hotfix_name,omitempty"`
	// Virtual foreign key for mapping entity.Repo
	RepoFullName string `json:"repo_full_name,omitempty"`

	BasedReleaseVersion string `json:"based_release_version,omitempty"`
	BasedCommitSha      string `json:"based_commit_sha,omitempty"`
	Branch              string `json:"branch,omitempty"`

	// Virtual foreign keys for mapping entity.Issue
	IssueIDs string `json:"issue_ids,omitempty"`
	// Virtual foreign keys for mapping entity.PullRequest
	MasterPrIDs  string `json:"master_pr_ids,omitempty"`
	BranchPrIDs  string `json:"branch_pr_ids,omitempty"`
	AllPRsPushed bool   `json:"all_prs_pushed,omitempty"`

	BuildID      string `json:"build_id,omitempty"`
	BuildFeature string `json:"build_feature,omitempty"`
}

// DB-Table
func (HotfixReleaseInfo) TableName() string {
	return "hotfix_release_info"
}

func (release HotfixReleaseInfo) ExtractIssueIds() []string {
	return strings.Split(release.IssueIDs, ",")
}

func (release HotfixReleaseInfo) ExtractMasterPrIds() []string {
	return strings.Split(release.MasterPrIDs, ",")
}

func (release HotfixReleaseInfo) ExtractBranchPrIds() []string {
	return strings.Split(release.BranchPrIDs, ",")
}

type HotfixReleaseInfoOptions struct {
	HotfixName  string   `json:"hotfix_name,omitempty" form:"hotfix_name"`
	HotfixNames []string `json:"hotfix_names,omitempty" form:"hotfix_names"`

	ListOption
}

type HotfixReleaseEntities []HotfixReleaseInfo

func (entities HotfixReleaseEntities) ExtractIssueIds() []string {
	issueIds := make([]string, 0)
	for _, entity := range entities {
		issueIds = append(issueIds, entity.ExtractIssueIds()...)
	}

	return issueIds
}

func (entities HotfixReleaseEntities) ExtractAllPrIds() []string {
	prIds := make([]string, 0)
	prIds = append(prIds, entities.ExtractMasterPrIds()...)
	prIds = append(prIds, entities.ExtractBranchPrIds()...)

	return prIds
}

func (entities HotfixReleaseEntities) ExtractMasterPrIds() []string {
	prIds := make([]string, 0)
	for _, entity := range entities {
		prIds = append(prIds, entity.ExtractMasterPrIds()...)
	}

	return prIds
}

func (entities HotfixReleaseEntities) ExtractBranchPrIds() []string {
	prIds := make([]string, 0)
	for _, entity := range entities {
		prIds = append(prIds, entity.ExtractBranchPrIds()...)
	}

	return prIds
}
