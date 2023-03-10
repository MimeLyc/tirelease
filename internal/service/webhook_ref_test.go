package service

import (
	"testing"
	"tirelease/commons/git"
	"tirelease/internal/store"
	"tirelease/utils/configs"

	"github.com/stretchr/testify/assert"
)

func TestValidateRefTestValidateRef(t *testing.T) {
	ref := "heads/release-6.1"
	refType := "branch"

	result := validateRef(ref, refType)
	assert.Equal(t, true, result)

	ref = "release-6.1-20221111"
	refType = "branch"

	result = validateRef(ref, refType)
	assert.Equal(t, false, result)

	ref = "release-6.1"
	refType = "tag"

	result = validateRef(ref, refType)
	assert.Equal(t, false, result)
}

func TestRefreshSprintTestRefreshSprint(t *testing.T) {
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	config := configs.NewConfig(
		"../../"+configs.TestConfig,
		"../../"+configs.TestSecretConfig)
	store.NewStore(config)

	major, minor := 6, 0
	owner, repo := "pingcap", "tidb"

	err := refreshSprint(major, minor, owner, repo)
	assert.Nil(t, err)
}
