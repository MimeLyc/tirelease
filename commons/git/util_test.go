package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseIssueUrl(t *testing.T) {
	url := "https://github.com/pingcap/tidb/issues/3333"

	org, repo, issue, err := ParseIssueUrl(url)
	assert.Nil(t, err)
	assert.Equal(t, "pingcap", org)
	assert.Equal(t, "tidb", repo)
	assert.Equal(t, 3333, issue)

	url = "https://github.com/pingcap/tidb/issues/"

	org, repo, issue, err = ParseIssueUrl(url)
	assert.NotNil(t, err)

}
