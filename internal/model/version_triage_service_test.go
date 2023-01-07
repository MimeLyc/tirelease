package model

import (
	"fmt"
	"testing"
	"tirelease/commons/configs"
	"tirelease/commons/database"
	"tirelease/internal/entity"
	"tirelease/internal/repository"

	"github.com/stretchr/testify/assert"
)

func TestGetRelatedPrs(t *testing.T) {
	configs.LoadConfig("../../config.yaml")
	config := configs.Config
	database.Connect(config)

	issueID := "I_kwDODAH3lM5Ly9_A"
	releaseBranch := "release-6.1"
	prs, err := SelectRelatedPrs(releaseBranch, issueID)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(prs))

	issueID = "I_kwDOCfvVlc5KD7UP"
	releaseBranch = "release-6.1"
	prs, err = SelectRelatedPrs(releaseBranch, issueID)
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, len(prs))

	issueID = "I_kwDODAH3lM5LXlGJ"
	releaseBranch = "release-6.1"
	prs, err = SelectRelatedPrs(releaseBranch, issueID)
	assert.Equal(t, nil, err)
	assert.Equal(t, 2, len(prs))
}

func TestExtractIssueIsFromTriages(t *testing.T) {
	configs.LoadConfig("../../config.yaml")
	config := configs.Config
	database.Connect(config)

	triages, err := repository.SelectVersionTriage(
		&entity.VersionTriageOption{
			VersionName: "6.5.0",
		},
	)

	assert.Nil(t, err)
	ids := extractIssueIDsFromTriage(*triages)
	fmt.Print(ids)
}
