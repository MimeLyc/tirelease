package model

import (
	"testing"
	"tirelease/commons/git"
	"tirelease/internal/entity"
	"tirelease/internal/store"
	"tirelease/utils/configs"

	"github.com/stretchr/testify/assert"
)

func TestApprovePr(t *testing.T) {
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)

	owner := "PingCAP-QE"
	repo := "tirelease"
	prNumber := 136

	git.Client.RemoveLabel(owner, repo, prNumber, git.CherryPickLabel)

	state, _ := NewPickTriageState(ParseFromEntityPickTriage(entity.VersionTriageResultLater))
	pickTriage := pickTriageStateContext{
		State: state,
		Version: &ReleaseVersion{
			ReleaseVersion: &entity.ReleaseVersion{
				Status: entity.ReleaseVersionStatusUpcoming,
			},
		},
		Prs: []PullRequest{
			{
				PullRequest: &entity.PullRequest{
					Owner:  owner,
					Repo:   repo,
					Number: prNumber,
				},
			},
		},
	}

	pickTriage.Trans(ParseFromEntityPickTriage(entity.VersionTriageResultAccept))
	pr, _, _ := git.Client.GetPullRequestByNumber(owner, repo, prNumber)
	labelsString := ""
	for _, label := range pr.Labels {
		labelsString = labelsString + *label.Name
	}

	assert.Contains(t, labelsString, git.CherryPickLabel)
	assert.Equal(t, pickTriage.State.getStateText(), ParseFromEntityPickTriage(entity.VersionTriageResultAccept))
}

func TestLeaveApprove(t *testing.T) {
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)

	owner := "PingCAP-QE"
	repo := "tirelease"
	prNumber := 136

	git.Client.AddLabel(owner, repo, prNumber, git.CherryPickLabel)

	state, _ := NewPickTriageState(ParseFromEntityPickTriage(entity.VersionTriageResultAccept))
	pickTriage := pickTriageStateContext{
		State: state,

		Version: &ReleaseVersion{
			ReleaseVersion: &entity.ReleaseVersion{
				Status: entity.ReleaseVersionStatusUpcoming,
			},
		},
		Prs: []PullRequest{
			{
				PullRequest: &entity.PullRequest{
					Owner:  owner,
					Repo:   repo,
					Number: prNumber,
				},
			},
		},
	}

	pickTriage.Trans(ParseFromEntityPickTriage(entity.VersionTriageResultLater))

	pr, _, _ := git.Client.GetPullRequestByNumber(owner, repo, prNumber)
	labelsString := ""
	for _, label := range pr.Labels {
		labelsString = labelsString + *label.Name
	}

	assert.NotContains(t, labelsString, git.CherryPickLabel)
	assert.Equal(t, ParseFromEntityPickTriage(entity.VersionTriageResultLater), pickTriage.State.getStateText())

}

func TestWontFixPr(t *testing.T) {
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	config := configs.NewConfig(
		"../../"+configs.TestConfig,
		"../../"+configs.TestSecretConfig)
	store.NewStore(config)

	owner := "PingCAP-QE"
	repo := "tirelease"
	prNumber := 136

	pr, _, _ := git.Client.GetPullRequestByNumber(owner, repo, prNumber)

	state, _ := NewPickTriageState(ParseFromEntityPickTriage(entity.VersionTriageResultLater))
	pickTriage := pickTriageStateContext{
		State: state,
		Version: &ReleaseVersion{
			ReleaseVersion: &entity.ReleaseVersion{
				Status: entity.ReleaseVersionStatusUpcoming,
			},
		},
		Issue: &entity.Issue{
			IssueID: "test",
		},
		Prs: []PullRequest{
			{
				PullRequest: &entity.PullRequest{
					PullRequestID: *pr.NodeID,
					Owner:         owner,
					Repo:          repo,
					Number:        prNumber,
				},
			},
		},
	}

	pickTriage.Trans(ParseFromEntityPickTriage(entity.VersionTriageResultWontFix))
}

func TestPickAccept(t *testing.T) {
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	config := configs.NewConfig(
		"../../"+configs.TestConfig,
		"../../"+configs.TestSecretConfig)
	store.NewStore(config)

	issueVersionTriage, err := SelectActiveIssueVersionTriage("5.4.3", "I_kwDOCfvVlc5S4W-q")
	assert.Nil(t, err)
	issueVersionTriage.Version.Status = entity.ReleaseVersionStatusUpcoming
	issueVersionTriage.PickTriage.State.StateText = StateText(entity.VersionTriageResultAcceptFrozen)
	issueVersionTriage.TriagePickStatus(entity.VersionTriageResultAccept)
	triage := issueVersionTriage.MapToEntity()
	assert.Equal(t, entity.VersionTriageResultAccept, triage.TriageResult)
}
