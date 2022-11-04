package model

import (
	"testing"
	"tirelease/internal/entity"

	"github.com/stretchr/testify/assert"
)

func TestTransBlock(t *testing.T) {
	context, err := NewBlockTriageStateContext(
		ParseFromEntityBlockTriage(entity.BlockVersionReleaseResultNoneBlock),
		&entity.Issue{},
		&ReleaseVersion{},
		nil,
	)

	assert.Nil(t, err)

	context.Trans(ParseFromEntityBlockTriage(entity.BlockVersionReleaseResultBlock))

	assert.Equal(t, context.State.getStateText(), ParseFromEntityBlockTriage(entity.BlockVersionReleaseResultBlock))

}
