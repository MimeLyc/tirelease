package model

import (
	"testing"
	"tirelease/commons/configs"
	"tirelease/commons/database"
	"tirelease/commons/git"
	"tirelease/internal/entity"
	"tirelease/internal/repository"

	"github.com/stretchr/testify/assert"
)

func TestChangeStateText(t *testing.T) {
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	configs.LoadConfig("../../config.yaml")
	config := configs.Config
	database.Connect(config)

	versionEntity, err := repository.SelectReleaseVersionLatest(
		&entity.ReleaseVersionOption{
			Name: "6.1.1",
		},
	)
	assert.Nil(t, err)
	version := Parse2ReleaseVersion(*versionEntity)
	assert.Equal(t, &version.Version, &version.versionStateContext.Version)
	version.ReleaseVersion.Status = entity.ReleaseVersionStatusPlanned
	assert.Equal(t, version.versionStateContext.Version.Status, entity.ReleaseVersionStatusPlanned)
}

func TestGetTransition(t *testing.T) {
	state, err := NewVersionState(StateText(entity.ReleaseVersionStatusPlanned))
	assert.Nil(t, err)

	meta := StateTransitionMeta{
		FromState: StateText(entity.ReleaseVersionStatusPlanned),
		ToState:   StateText(entity.ReleaseVersionStatusUpcoming),
	}
	transition := state.getTransition(meta)
	assert.Equal(t, VersionState2Upcoming{}, transition)

	meta = StateTransitionMeta{
		FromState: StateText(entity.ReleaseVersionStatusPlanned),
		ToState:   StateText(entity.ReleaseVersionStatusFrozen),
	}
	transition = state.getTransition(meta)
	assert.Equal(t, VersionState2Frozen{}, transition)
}
