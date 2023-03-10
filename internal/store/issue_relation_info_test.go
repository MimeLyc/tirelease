package store

import (
	"testing"

	"tirelease/commons/git"
	"tirelease/internal/entity"
	"tirelease/utils/configs"

	"github.com/stretchr/testify/assert"
)

func TestSelectIssueRelationInfoByJoin(t *testing.T) {
	// Init
	config := configs.NewConfig(
		"../../"+configs.TestConfig,
		"../../"+configs.TestSecretConfig)
	NewStore(config)

	// Option
	option := &entity.IssueRelationInfoOption{
		IssueOption: entity.IssueOption{
			IssueID: git.TestIssueNodeID,
		},
		AffectVersion: "5.4",
	}
	issueRelationInfoJoin, err := SelectIssueRelationInfoByJoin(option)

	// Assert
	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, len(*issueRelationInfoJoin) > 0)

}

func TestSelectUnpickedIssueRelationInfo(t *testing.T) {
	// Init
	config := configs.NewConfig(
		"../../"+configs.TestConfig,
		"../../"+configs.TestSecretConfig)
	NewStore(config)

	// Option
	option := &entity.IssueRelationInfoOption{
		AffectVersion: "5.4",
	}
	issueRelationInfoJoin, err := SelectNeedTriageIssueRelationInfo(option)

	// Assert
	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, len(*issueRelationInfoJoin) > 0)

}
