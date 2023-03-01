package service

import (
	"testing"

	"tirelease/commons/database"
	"tirelease/commons/git"
	"tirelease/internal/entity"
	"tirelease/internal/model"

	"github.com/google/go-github/v41/github"
	"github.com/stretchr/testify/assert"
)

func TestCronRefreshPullRequestV4(t *testing.T) {
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
	params := &RefreshPullRequestParams{
		Repos:       &repos,
		BeforeHours: -2,
		Batch:       20,
		Total:       500,
	}

	// detail
	err := CronRefreshPullRequestV4(params)
	assert.Equal(t, true, err == nil)
}

func TestCronMergeRetryPullRequestV3(t *testing.T) {
	t.Skip()
	// init
	git.Connect(git.TestToken)
	database.Connect(generateConfig())

	// detail
	err := CronMergeRetryPullRequestV3()
	assert.Equal(t, true, err == nil)
}

func TestWebhookRefreshPullRequestV3(t *testing.T) {
	// init
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	database.Connect(generateConfig())

	// detail
	pr, _, err := git.Client.GetPullRequestByNumber("pingcap", "tidb", 38066)
	assert.Equal(t, true, err == nil)
	err = WebhookRefreshPullRequestV3(pr)
	assert.Equal(t, true, err == nil)
}

func TestCronRefreshPullRequestV42(t *testing.T) {
	t.Skip()
	// init
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	database.Connect(generateConfig())
	repo := &entity.Repo{
		Owner: "pingcap",
		Repo:  "tiflow",
	}
	repos := []entity.Repo{*repo}
	params := &RefreshPullRequestParams{
		Repos:       &repos,
		BeforeHours: -4380,
		Batch:       20,
		Total:       3000,
	}

	// detail
	err := CronRefreshPullRequestV4(params)
	assert.Equal(t, true, err == nil)
}

func TestWebHookRefreshPullRequestRefIssue(t *testing.T) {
	t.Skip()
	// init
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	database.Connect(generateConfig())

	pr := &github.PullRequest{
		NodeID: &git.TestPullRequestNodeID3,
	}
	err := WebHookRefreshPullRequestRefIssue(pr)
	assert.Equal(t, true, err == nil)
}

func TestCheckTriageStatus(t *testing.T) {
	// t.Skip()
	database.Connect(generateConfig())
	issues, err := model.IssueCmd{
		IssueOption: &entity.IssueOption{
			Number: 6851,
			Owner:  "pingcap",
			Repo:   "tiflash",
		},
		AffectOption: &entity.IssueAffectOption{
			AffectVersion: "6.5",
			AffectResult:  entity.AffectResultResultYes,
		},
		TriageBuildCommand: &model.TriageBuildCommand{
			WithTriages: true,
		},
	}.BuildArray()

	allApproved, err := checkTriageStatus("6.5", issues)
	assert.Equal(t, true, allApproved)
	assert.Equal(t, nil, err)

}

func TestAutoRefreshPrApprovedLabel(t *testing.T) {
	t.Skip()
	// init
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	database.Connect(generateConfig())

	pr, _, _ := git.Client.GetPullRequestByNumber("pingcap", "tiflash", 4970)
	AutoRefreshPrApprovedLabel(pr)
}

func TestRefreshPrIssueRefByPrContent(t *testing.T) {
	t.Skip()

	// init
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	database.Connect(generateConfig())

	prV4, err := git.ClientV4.GetPullRequestByID(git.TestPullRequestNodeID3)
	assert.Equal(t, nil, err)
	refreshPrIssueRefByPrContent(prV4)

}
