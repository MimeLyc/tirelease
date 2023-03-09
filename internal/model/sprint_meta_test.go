package model

import (
	"testing"
	"time"
	"tirelease/commons/git"
	"tirelease/internal/entity"
	"tirelease/internal/store"
	"tirelease/utils/configs"

	"github.com/stretchr/testify/assert"
)

func TestSprintStartTime(t *testing.T) {
	git.ConnectV4(git.TestToken)
	config := configs.NewConfig(
		"../../"+configs.TestConfig,
		"../../"+configs.TestSecretConfig)
	store.NewStore(config)

	sprint := SprintMeta{
		entity.SprintMeta{
			Major: 6,
			Minor: 1,
			Repo: entity.Repo{
				Owner: "pingcap",
				Repo:  "tidb",
			},
		},
	}

	startTime, err := CalculateStartTimeOfSprint(sprint.Major, sprint.Minor, sprint.Repo)
	assert.Nil(t, err)
	assert.Equal(t, time.Time(time.Date(2022, time.March, 17, 12, 42, 31, 0, time.UTC)), *startTime)
}
