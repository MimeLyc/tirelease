package model

import (
	"testing"
	"tirelease/commons/configs"
	"tirelease/commons/database"
	"tirelease/commons/git"
	"tirelease/internal/entity"

	"github.com/stretchr/testify/assert"
)

func TestIssueBuilder(t *testing.T) {
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	configs.LoadConfig("../../config.yaml")
	config := configs.Config
	database.Connect(config)

	builder := IssueBuilder{}
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
			NeedTriages: true,
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
