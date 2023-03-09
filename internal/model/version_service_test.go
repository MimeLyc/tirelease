package model

import (
	"fmt"
	"testing"
	"tirelease/commons/configs"
	"tirelease/commons/database"
	"tirelease/internal/entity"
	"tirelease/internal/store"

	"github.com/stretchr/testify/assert"
)

func TestSelectHistoryIssueTriage(t *testing.T) {
	configs.LoadConfig("../../config.yaml")
	config := configs.Config
	database.Connect(config)

	versionEntity, err := store.SelectReleaseVersionLatest(
		&entity.ReleaseVersionOption{
			Name: "6.5.0",
		},
	)
	assert.Nil(t, err)

	version := Parse2ReleaseVersion(*versionEntity)

	triages, err := version.SelectHistoryIssueTriages()
	fmt.Print(triages)
}
