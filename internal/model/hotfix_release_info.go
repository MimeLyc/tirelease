package model

import (
	"fmt"
	"strings"
	"time"
	"tirelease/commons/git"
	"tirelease/internal/entity"
)

type HotfixReleaseInfo struct {
	entity.HotfixReleaseInfo
	Issues    []Issue       `json:"issues,omitempty"`
	MasterPrs []PullRequest `json:"master_prs,omitempty"`
	BranchPrs []PullRequest `json:"branch_prs,omitempty"`
	Assignee  *User         `json:"assignee,omitempty"`
}

// Find or create hotfix release branch
func (release HotfixReleaseInfo) FetchHotfixBranch(baseVersion string) (string, error) {
	owner := strings.Split(release.RepoFullName, "/")[0]
	repo := strings.Split(release.RepoFullName, "/")[1]
	tag := release.BasedReleaseVersion

	commitId, err := GetCommitByTag(owner, repo, tag)
	branches, err := GetBranchesByCommit(owner, repo, commitId)
	if err != nil {
		return "", err
	}
	branches = filterTargetHotfixBranch(baseVersion, branches)

	if len(branches) > 1 {
		return "", fmt.Errorf("found multiple hotfix branches: %v", branches)
	} else if len(branches) == 1 {
		return branches[0], nil
	}

	// Create hotfix Branch
	branchName := fmt.Sprintf(
		git.HotfixBranch,
		ExtractVersionMinorName(baseVersion),
		time.Now().Format("20060102"),
		baseVersion,
	)
	err = CreateBranchByCommit(owner, repo, branchName, commitId)
	if err != nil {
		return "", err
	}
	return branchName, nil
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

func filterTargetHotfixBranch(baseVersion string, branches []string) []string {
	minorName := ExtractVersionMinorName(baseVersion)

	result := make([]string, 0)
	for _, branch := range branches {
		if !git.IsHotfixBranch(branch) {
			continue
		}

		if strings.HasPrefix(branch, fmt.Sprintf(git.HotfixBranchPrefix, minorName)) {
			result = append(result, branch)
		}
	}
	return result
}
