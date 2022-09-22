package model

import (
	"testing"
	"tirelease/commons/git"

	"github.com/stretchr/testify/assert"
)

func TestGetCheckoutCommit(t *testing.T) {
	git.ConnectV4(git.TestToken)
	owner := "pingcap"
	repo := "tidb"
	ref := "release-6.3"
	commit, error := GetCheckoutCommitOfRef(owner, repo, ref, git.RefTypeBranch)
	assert.Nil(t, error)
	assert.Equal(t, "974b5784adbbd47d14659916d47dd986effa7b4e", commit.Oid)
}
