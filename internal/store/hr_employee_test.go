package store

import (
	"testing"
	"tirelease/utils/configs"

	"github.com/stretchr/testify/assert"
)

func TestSelectAllHrEmployees(t *testing.T) {
	config := configs.NewConfig(
		"../../"+configs.TestConfig,
		"../../"+configs.TestSecretConfig)
	NewStore(config)

	hrEmployees, err := SelectAllHrEmployee()
	assert.Nil(t, err)
	assert.NotNil(t, hrEmployees)
}
