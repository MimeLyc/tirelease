package model

import (
	"testing"
	"tirelease/commons/configs"
	"tirelease/commons/database"
	"tirelease/commons/git"
	"tirelease/internal/entity"
	"tirelease/internal/store"

	"github.com/stretchr/testify/assert"
)

func TestNewIssueVersionTriageNoHistory(t *testing.T) {
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	configs.LoadConfig("../../config.yaml")
	config := configs.Config
	database.Connect(config)

	issueVersionTriage, err := SelectActiveIssueVersionTriage("6.1", "MDU6SXNzdWU4MTc4NzUyNDQ=")
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(0), issueVersionTriage.ID)
	assert.Equal(t, "6.1.1", issueVersionTriage.Version.Name)
	assert.Equal(t, 1498, issueVersionTriage.Issue.Number)
	assert.Equal(t, 0, len(issueVersionTriage.RelatedPrs))
	assert.Equal(t, EmptyStateText(), issueVersionTriage.PickTriage.State.StateText)
	assert.Equal(t, ParseFromEntityBlockTriage(entity.BlockVersionReleaseResultBlock), issueVersionTriage.BlockTriage.State.StateText)
}

func TestNewIssueVersionTriageWithHistory(t *testing.T) {
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	configs.LoadConfig("../../config.yaml")
	config := configs.Config
	database.Connect(config)

	issueVersionTriage, err := SelectActiveIssueVersionTriage("6.1", "I_kwDOAoCpQc5OQby4")
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(1260002), issueVersionTriage.ID)
	assert.Equal(t, "6.1.1", issueVersionTriage.Version.Name)
	assert.Equal(t, 36426, issueVersionTriage.Issue.Number)
	assert.Equal(t, 1, len(issueVersionTriage.RelatedPrs))
	assert.Equal(t, ParseFromEntityPickTriage(entity.VersionTriageResultWontFix), issueVersionTriage.PickTriage.State.StateText)
	assert.Equal(t, EmptyStateText(), issueVersionTriage.BlockTriage.State.StateText)
}

func TestMapToEntityWithHistory(t *testing.T) {
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	configs.LoadConfig("../../config.yaml")
	config := configs.Config
	database.Connect(config)

	rawVersionTriage, err := store.SelectVersionTriageUnique(
		&entity.VersionTriageOption{
			ID: int64(1260002),
		},
	)
	assert.Equal(t, nil, err)
	issueVersionTriage, err := SelectActiveIssueVersionTriage("6.1", "I_kwDOAoCpQc5OQby4")
	assert.Equal(t, nil, err)
	parsedVersionTriage := issueVersionTriage.MapToEntity()

	assert.Equal(t, rawVersionTriage.VersionName, parsedVersionTriage.VersionName)
	assert.Equal(t, rawVersionTriage.ID, parsedVersionTriage.ID)
	assert.Equal(t, rawVersionTriage.IssueID, parsedVersionTriage.IssueID)
	assert.Equal(t, rawVersionTriage.TriageResult, parsedVersionTriage.TriageResult)
	assert.Equal(t, rawVersionTriage.BlockVersionRelease, parsedVersionTriage.BlockVersionRelease)
}

func TestEmptyState(t *testing.T) {
	test := StateText("")
	assert.Equal(t, EmptyStateText(), test)
}
