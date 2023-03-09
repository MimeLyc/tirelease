package configs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	//t.Skip()
	config := NewConfig("../../"+TestConfig, "../../"+TestSecretConfig)
	assert.Equal(t, true, config.GitHubAccessToken != "")
}
