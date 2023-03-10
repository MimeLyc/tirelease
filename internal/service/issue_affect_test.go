package service

import (
	"testing"
	"tirelease/internal/store"
	"tirelease/utils/configs"

	"tirelease/commons/git"
	"tirelease/internal/entity"

	"github.com/stretchr/testify/assert"
)

func TestComposeIssueAffectWithIssueID(t *testing.T) {
	// Init
	config := configs.NewConfig(
		"../../"+configs.TestConfig,
		"../../"+configs.TestSecretConfig)
	store.NewStore(config)

	// Test
	issueAffects, err := ComposeIssueAffectWithIssueID(git.TestIssueNodeID, nil)
	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, len(*issueAffects) > 0)
}

func TestCreateOrUpdateIssueAffect(t *testing.T) {
	t.Skip()
	// init
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	config := configs.NewConfig(
		"../../"+configs.TestConfig,
		"../../"+configs.TestSecretConfig)
	store.NewStore(config)

	// Test
	issueAffect := &entity.IssueAffect{
		IssueID:       git.TestIssueNodeID2,
		AffectVersion: "6.0",
		AffectResult:  entity.AffectResultResultYes,
	}
	err := CreateOrUpdateIssueAffect(issueAffect)
	assert.Equal(t, true, err == nil)
}
