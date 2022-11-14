package model

import (
	"fmt"
	"testing"
	"tirelease/commons/configs"
	"tirelease/commons/database"
	"tirelease/commons/git"

	"github.com/stretchr/testify/assert"
)

func TestSelectIssueTriage(t *testing.T) {
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	configs.LoadConfig("../../config.yaml")
	config := configs.Config
	database.Connect(config)

	issueId := "I_kwDOAuklds5RIZht"

	issueTriage, err := SelectIssueTriage(issueId)

	assert.Nil(t, err)

	fmt.Printf("%v", issueTriage)
}
