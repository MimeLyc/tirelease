package store

import (
	"testing"
	"tirelease/utils/configs"

	"github.com/stretchr/testify/assert"
)

func TestBatchSelectEmployeesByGhLogins(t *testing.T) {
	config := configs.NewConfig(
		"../../"+configs.TestConfig,
		"../../"+configs.TestSecretConfig)
	NewStore(config)

	githubLogins := []string{"MimeLyc", "VelocityLight"}
	employees, err := BatchSelectEmployeesByGhLogins(githubLogins)
	assert.Nil(t, err)
	assert.NotNil(t, employees)
}
