package store

import (
	"testing"
	"tirelease/utils/configs"

	"github.com/stretchr/testify/assert"
	"tirelease/commons/git"
	"tirelease/internal/entity"
)

func TestSelectIssueRaw(t *testing.T) {
	// Init
	config := configs.NewConfig(
		"../../"+configs.TestConfig,
		"../../"+configs.TestSecretConfig)
	NewStore(config)

	// Select
	var option = &entity.IssueOption{
		IssueIDs: []string{git.TestIssueNodeID, git.TestIssueNodeID2},

		ListOption: entity.ListOption{
			Page:    1,
			PerPage: 10,

			OrderBy: "id",
			Order:   "desc",
		},
	}
	issues, err := SelectIssue(option)
	count, _ := CountIssue(option)
	// Assert
	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, len(*issues) > 0)
	assert.Equal(t, true, count > 0)
}
