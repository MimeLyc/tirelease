package store

import (
	"testing"

	"tirelease/commons/git"
	"tirelease/internal/entity"
	"tirelease/utils/configs"

	"github.com/stretchr/testify/assert"
)

func TestPullRequest(t *testing.T) {
	// Init
	config := configs.NewConfig(
		"../../"+configs.TestConfig,
		"../../"+configs.TestSecretConfig)
	NewStore(config)

	// Select
	var option = &entity.PullRequestOption{
		PullRequestID: git.TestPullRequestNodeID,

		ListOption: entity.ListOption{
			Page:    1,
			PerPage: 10,

			OrderBy: "id",
			Order:   "desc",
		},
	}
	prs, err := SelectPullRequest(option)
	// Assert
	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, len(*prs) > 0)
}
