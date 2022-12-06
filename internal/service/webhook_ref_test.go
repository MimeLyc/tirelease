package service

import (
	"testing"
	"tirelease/commons/database"
	"tirelease/commons/git"

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
	database.Connect(generateConfig())

	major, minor := 6, 0
	owner, repo := "pingcap", "tidb"

	err := refreshSprint(major, minor, owner, repo)
	assert.Nil(t, err)
}
