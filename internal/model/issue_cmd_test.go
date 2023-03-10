package model

import (
	"testing"
	"tirelease/commons/git"
	"tirelease/internal/entity"
	"tirelease/internal/store"
	"tirelease/utils/configs"

	"github.com/stretchr/testify/assert"
)

func TestIssueBuilder(t *testing.T) {
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	config := configs.NewConfig(
		"../../"+configs.TestConfig,
		"../../"+configs.TestSecretConfig)
	store.NewStore(config)

	builder := IssueCmd{}
	builder = builder.Option(
		&entity.IssueOption{
			Repo: "tidb",
		},
		&entity.IssueAffectOption{
			AffectVersion: "6.5",
			AffectResult:  "Yes",
		},
	).Command(
		&TriageBuildCommand{
			WithTriages: true,
		},
	)

	issues, err := builder.BuildArray()
	assert.Nil(t, err)
	for _, issue := range issues {
		for _, triage := range issue.VersionTriages {
			assert.Equal(t, triage.IssueID, issue.IssueID)
		}
	}
}
