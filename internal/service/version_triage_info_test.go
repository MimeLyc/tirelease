package service

import (
	"testing"

	"tirelease/commons/git"
	"tirelease/internal/entity"
	"tirelease/internal/store"
	"tirelease/utils/configs"

	"github.com/stretchr/testify/assert"
)

func TestCreateOrUpdateVersionTriageInfo(t *testing.T) {
	t.Skip()
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	config := configs.NewConfig(
		"../../"+configs.TestConfig,
		"../../"+configs.TestSecretConfig)
	store.NewStore(config)

	versionTriage := &entity.VersionTriage{
		VersionName:  "6.0",
		IssueID:      git.TestIssueNodeID2,
		TriageResult: entity.VersionTriageResultUnKnown,
	}
	info, err := SaveVersionTriageInfo(versionTriage)

	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, info != nil)
	assert.Equal(t, true, info.IsAccept)
}

func TestComposeVersionTriageUpcomingList(t *testing.T) {
	t.Skip()
	config := configs.NewConfig(
		"../../"+configs.TestConfig,
		"../../"+configs.TestSecretConfig)
	store.NewStore(config)

	versionTriages, err := ComposeVersionTriageUpcomingList("5.0.7")
	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, len(versionTriages) > 0)
}

func TestChangePrApprovedLabel(t *testing.T) {
	t.Skip()
	config := configs.NewConfig(
		"../../"+configs.TestConfig,
		"../../"+configs.TestSecretConfig)
	store.NewStore(config)
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)

	pr, _, err := git.Client.GetPullRequestByNumber("PingCAP-QE", "tirelease", 111)
	assert.Equal(t, true, err == nil)
	err = ChangePrApprovedLabel(*pr.NodeID, false, true)
	assert.Equal(t, true, err == nil)

	err = ChangePrApprovedLabel(*pr.NodeID, true, false)
	assert.Equal(t, true, err == nil)
}
