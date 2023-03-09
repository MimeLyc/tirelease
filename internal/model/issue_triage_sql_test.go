package model

import (
	"fmt"
	"testing"
	"tirelease/commons/git"
	"tirelease/internal/store"
	"tirelease/utils/configs"

	"github.com/stretchr/testify/assert"
)

func TestSelectIssueTriage(t *testing.T) {
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	config := configs.NewConfig(
		"../../"+configs.TestConfig,
		"../../"+configs.TestSecretConfig)
	store.NewStore(config)

	issueId := "I_kwDOAuklds5RIZht"

	issueTriage, err := SelectIssueTriage(issueId)

	assert.Nil(t, err)

	fmt.Printf("%v", issueTriage)
}
