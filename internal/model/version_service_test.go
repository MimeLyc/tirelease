package model

import (
	"fmt"
	"testing"
	"tirelease/internal/entity"
	"tirelease/internal/store"
	"tirelease/utils/configs"

	"github.com/stretchr/testify/assert"
)

func TestSelectHistoryIssueTriage(t *testing.T) {
	config := configs.NewConfig(
		"../../"+configs.TestConfig,
		"../../"+configs.TestSecretConfig)
	store.NewStore(config)

	versionEntity, err := store.SelectReleaseVersionLatest(
		&entity.ReleaseVersionOption{
			Name: "6.5.0",
		},
	)
	assert.Nil(t, err)

	version := Parse2ReleaseVersion(*versionEntity)

	triages, err := version.SelectHistoryIssueTriages()
	fmt.Print(triages)
	assert.NoError(t, err)
}
