package service

import (
	"testing"
	"time"

	"tirelease/commons/git"
	"tirelease/internal/entity"
	"tirelease/internal/store"
	"tirelease/utils/configs"

	"github.com/stretchr/testify/assert"
)

func TestGetPullRequestByNumberFromV3(t *testing.T) {
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)

	pr, err := GetPRByNumberFromV3(git.TestOwner, git.TestRepo, git.TestPullRequestId)

	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, pr != nil)
}

func TestGetPullRequestRefIssuesByRegexFromV4(t *testing.T) {
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)

	pr, err := git.ClientV4.GetPullRequestByID(git.TestPullRequestNodeID)
	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, pr != nil)

	issueNumbers, err := GetPullRequestRefIssuesByRegexFromV4(pr)
	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, len(issueNumbers) > 0)
}

func TestRegexReferenceNumbers(t *testing.T) {
	s := "close #1"
	issueNumbers, err := RegexReferenceNumbers(s)
	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, len(issueNumbers) == 1)

	s = "close #10, #100, #1000"
	issueNumbers, err = RegexReferenceNumbers(s)
	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, len(issueNumbers) == 3)
}

// script for refresh pull reqeust info
func TestRefreshPullRequestInfo(t *testing.T) {
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	config := configs.NewConfig(
		"../../"+configs.TestConfig,
		"../../"+configs.TestSecretConfig)
	store.NewStore(config)
	isMerged := true

	prsToRefresh := make([]entity.PullRequest, 0)
	// start time is the last sprint checkout time.
	startTime, _ := time.Parse("2006-01-02", "2022-09-20")
	// Select PRs on default branch
	prs, _ := store.SelectPullRequest(
		&entity.PullRequestOption{
			BaseBranch: "master",
			Merged:     &isMerged,
			MergeTime:  &startTime,
		},
	)
	for _, pr := range *prs {
		if pr.AuthorGhLogin == "" {
			prsToRefresh = append(prsToRefresh, pr)
		}
	}

	prs, _ = store.SelectPullRequest(
		&entity.PullRequestOption{
			BaseBranch: "main",
			Merged:     &isMerged,
			MergeTime:  &startTime,
		},
	)
	for _, pr := range *prs {
		if pr.AuthorGhLogin == "" {
			prsToRefresh = append(prsToRefresh, pr)
		}
	}

	// Select PRs on target sprint branch
	prs, _ = store.SelectPullRequest(
		&entity.PullRequestOption{
			BaseBranch: "release-6.5",
			Merged:     &isMerged,
			MergeTime:  &startTime,
		},
	)
	for _, pr := range *prs {
		if pr.AuthorGhLogin == "" {
			prsToRefresh = append(prsToRefresh, pr)
		}
	}

	for _, pr := range prsToRefresh {
		refreshedPr, _ := GetPRByNumberFromV3(pr.Owner, pr.Repo, pr.Number)
		store.CreateOrUpdatePullRequest(refreshedPr)
	}
}
