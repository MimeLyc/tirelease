package git

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHotfix(t *testing.T) {

	compile, err := regexp.Compile(HotfixBranchRegexP)
	assert.Nil(t, err)

	assert.True(t, compile.MatchString("release-6.1-20230301-v6.1.1"))
	assert.True(t, compile.MatchString("release-6.1-20230301"))
	assert.False(t, compile.MatchString("release-6.1.1-20230301"))
	assert.False(t, compile.MatchString("release-6.1"))
}
