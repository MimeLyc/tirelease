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
	ref := "v5.4.0"
	commit, error := GetCheckoutCommitOfRef(owner, repo, ref, git.RefTypeTag)
	assert.Nil(t, error)
	assert.Equal(t, "974b5784adbbd47d14659916d47dd986effa7b4e", commit.Oid)
}
