package service

import (
	"testing"
	"time"

	"tirelease/commons/database"
	"tirelease/commons/git"
	"tirelease/internal/entity"
	"tirelease/internal/model"
	"tirelease/internal/repository"

	"github.com/google/go-github/v41/github"
	"github.com/stretchr/testify/assert"
)

func TestGetIssueByNumberFromV3(t *testing.T) {
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)

	issue, err := GetIssueByNumberFromV3(git.TestOwner, git.TestRepo, git.TestIssueId)

	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, issue != nil)
}

func TestGetIssuesByTimeFromV3(t *testing.T) {
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)

	day, _ := time.ParseDuration("-24h")
	time := time.Now().Add(15 * day)
	issues, err := GetIssuesByTimeFromV3(git.TestOwner, git.TestRepo, &time)

	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, len(issues) > 0)
}

func TestGetIssuesByOptionV3(t *testing.T) {
	t.Skip()

	git.Connect(git.TestToken)

	newLabel := "affects-6.0"
	labels := []string{"severity/major", "type/bug"}
	option := github.IssueListByRepoOptions{
		State:  "open",
		Labels: labels,
	}
	repos := []entity.Repo{
		{
			Owner: "VelocityLight",
			Repo:  "tirelease",
		},
		{
			Owner: "pingcap",
			Repo:  "tidb",
		},
		{
			Owner: "pingcap",
			Repo:  "tiflash",
		},
		{
			Owner: "pingcap",
			Repo:  "tidb-binlog",
		},
		{
			Owner: "pingcap",
			Repo:  "br",
		},
		{
			Owner: "pingcap",
			Repo:  "tidb-tools",
		},
		{
			Owner: "pingcap",
			Repo:  "ticdc",
		},
		{
			Owner: "pingcap",
			Repo:  "dumpling",
		},
		{
			Owner: "tikv",
			Repo:  "tikv",
		},
		{
			Owner: "tikv",
			Repo:  "pd",
		},
		{
			Owner: "tikv",
			Repo:  "importer",
		},
	}

	for _, repo := range repos {
		issues, err := GetIssuesByOptionV3(repo.Owner, repo.Repo, &option)
		assert.Equal(t, true, err == nil)
		assert.Equal(t, true, len(issues) > 0)
		err = BatchLabelIssues(issues, newLabel)
		assert.Equal(t, true, err == nil)
	}
}

// script for refresh issue info
func TestRefreshIssueInfo(t *testing.T) {
	// t.Skip()

	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	database.Connect(generateConfig())

	//select issues affects version
	minorVersion := "6.4"
	affects, err := repository.SelectIssueAffect(
		&entity.IssueAffectOption{
			AffectVersion: minorVersion,
			AffectResult:  entity.AffectResultResultYes,
		},
	)
	assert.Nil(t, err)
	issueIds := model.ExtractIssueIDs(*affects)

	issuesToRefresh, err := repository.SelectIssue(
		&entity.IssueOption{
			IssueIDs: issueIds,
		},
	)
	assert.Nil(t, err)

	for _, issue := range *issuesToRefresh {
		refreshedIssue, err := GetIssueByNumberFromV3(issue.Owner, issue.Repo, issue.Number)
		assert.Nil(t, err)
		repository.CreateOrUpdateIssue(refreshedIssue)
	}
}

// script for refresh issue on master info
func TestRefreshMasterBugOfSpringInfo(t *testing.T) {
	// t.Skip()

	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	database.Connect(generateConfig())

	//select issues affects version
	major := 6
	minor := 3
	issuesToRefresh, err := model.SelectFixedBugsBeforeSprintCheckout(major, minor)
	assert.Nil(t, err)

	for _, issue := range issuesToRefresh {
		refreshedIssue, err := GetIssueByNumberFromV3(issue.Owner, issue.Repo, issue.Number)
		assert.Nil(t, err)
		repository.CreateOrUpdateIssue(refreshedIssue)
	}
}
