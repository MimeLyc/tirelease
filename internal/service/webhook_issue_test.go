package service

import (
	"testing"

	"tirelease/commons/database"
	"tirelease/commons/git"
	"tirelease/internal/entity"
	"tirelease/internal/repository"

	"github.com/stretchr/testify/assert"
)

func TestCronRefreshIssuesV4(t *testing.T) {
	t.Skip()
	// init
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	database.Connect(generateConfig())
	repo := &entity.Repo{
		Owner: git.TestOwner2,
		Repo:  git.TestRepo2,
	}
	repos := []entity.Repo{*repo}
	params := &RefreshIssueParams{
		Repos:       &repos,
		BeforeHours: -25,
		Batch:       20,
		Total:       500,
		IsHistory:   false,
		Order:       "DESC",
	}

	// detail
	err := CronRefreshIssuesV4(params)
	assert.Equal(t, true, err == nil)
}

func TestWebhookRefreshIssueV4(t *testing.T) {
	// init
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	database.Connect(generateConfig())

	// detail
	issue, _, err := git.Client.GetIssueByNumber(git.TestOwner2, git.TestRepo2, git.TestIssueId2)
	assert.Equal(t, true, err == nil)
	err = WebhookRefreshIssueV4(issue)
	assert.Equal(t, true, err == nil)
}

func TestCronRefreshIssuesV42(t *testing.T) {
	t.Skip()
	// init
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	database.Connect(generateConfig())
	repos, err := repository.SelectRepo(&entity.RepoOption{})
	if err != nil {
		return
	}
	releaseVersions, err := repository.SelectReleaseVersion(&entity.ReleaseVersionOption{})
	if err != nil {
		return
	}

	params := &RefreshIssueParams{
		Repos: repos,
		// BeforeHours:     -8760,
		BeforeHours:     -720,
		Batch:           10,
		Total:           20,
		IsHistory:       true,
		ReleaseVersions: releaseVersions,
		Order:           "DESC",
	}

	// detail
	err = CronRefreshIssuesV4(params)
	assert.Equal(t, true, err == nil)
}

func TestRefreshIssueField(t *testing.T) {
	t.Skip()
	// init
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	database.Connect(generateConfig())

	// detail
	option := &entity.IssueOption{
		ListOption: entity.ListOption{
			OrderBy: "id",
			Order:   "ASC",
		},
	}
	err := RefreshIssueField(option)
	assert.Equal(t, true, err == nil)
}
