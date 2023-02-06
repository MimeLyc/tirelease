package model

import (
	"testing"
	"tirelease/commons/configs"
	"tirelease/commons/database"
	"tirelease/commons/git"
	"tirelease/internal/entity"

	"github.com/stretchr/testify/assert"
)

func TestApprove(t *testing.T) {
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	configs.LoadConfig("../../config.yaml")
	config := configs.Config
	database.Connect(config)

	owner := "PingCAP-QE"
	repo := "tirelease"
	number := 135

	gitPr, err := git.ClientV4.GetPullRequestsByNumber(owner, repo, number)
	assert.Equal(t, nil, err)

	pr := PullRequest{
		PullRequest: &entity.PullRequest{
			PullRequestID: gitPr.ID.(string),
			Owner:         owner,
			Repo:          repo,
			Number:        number,
		},
	}

	err = pr.Approve()
	assert.Equal(t, nil, err)

	err = pr.UnApprove()
	assert.Equal(t, nil, err)

}

func TestBuildByIssues(t *testing.T) {
	configs.LoadConfig("../../config.yaml")
	config := configs.Config
	database.Connect(config)

	// Issue without related pr
	issueId := "MDU6SXNzdWU3NzIwOTEwOTM="
	prs, err := PullRequestCmd{
		ByRelatedIssue: true,
		IssueIds:       []string{issueId},
	}.Build()
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, len(prs))
}
