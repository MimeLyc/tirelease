package service

import (
	"testing"

	"tirelease/commons/git"
	"tirelease/internal/entity"
	"tirelease/internal/model"
	"tirelease/internal/store"
	"tirelease/utils/configs"

	"github.com/google/go-github/v41/github"
	"github.com/stretchr/testify/assert"
)

func TestCronRefreshPullRequestV4(t *testing.T) {
	t.Skip()
	// init
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	config := configs.NewConfig(
		"../../"+configs.TestConfig,
		"../../"+configs.TestSecretConfig)
	store.NewStore(config)
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
	config := configs.NewConfig(
		"../../"+configs.TestConfig,
		"../../"+configs.TestSecretConfig)
	store.NewStore(config)

	// detail
	err := CronMergeRetryPullRequestV3()
	assert.Equal(t, true, err == nil)
}

func TestWebhookRefreshPullRequestV3(t *testing.T) {
	// init
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	config := configs.NewConfig(
		"../../"+configs.TestConfig,
		"../../"+configs.TestSecretConfig)
	store.NewStore(config)

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
	config := configs.NewConfig(
		"../../"+configs.TestConfig,
		"../../"+configs.TestSecretConfig)
	store.NewStore(config)
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
	config := configs.NewConfig(
		"../../"+configs.TestConfig,
		"../../"+configs.TestSecretConfig)
	store.NewStore(config)

	pr := &github.PullRequest{
		NodeID: &git.TestPullRequestNodeID3,
	}
	err := WebHookRefreshPullRequestRefIssue(pr)
	assert.Equal(t, true, err == nil)
}

func TestCheckTriageStatus(t *testing.T) {
	// t.Skip()
	config := configs.NewConfig(
		"../../"+configs.TestConfig,
		"../../"+configs.TestSecretConfig)
	store.NewStore(config)
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
	config := configs.NewConfig(
		"../../"+configs.TestConfig,
		"../../"+configs.TestSecretConfig)
	store.NewStore(config)

	pr, _, _ := git.Client.GetPullRequestByNumber("pingcap", "tiflash", 4970)
	AutoRefreshPrApprovedLabel(pr)
}

func TestRefreshPrIssueRefByPrContent(t *testing.T) {
	t.Skip()

	// init
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	config := configs.NewConfig(
		"../../"+configs.TestConfig,
		"../../"+configs.TestSecretConfig)
	store.NewStore(config)

	prV4, err := git.ClientV4.GetPullRequestByID(git.TestPullRequestNodeID3)
	assert.Equal(t, nil, err)
	refreshPrIssueRefByPrContent(prV4)

}
