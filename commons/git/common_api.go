package git

import (
	"fmt"
	"strings"
)

// Distinguish issue or pull request method
func IsIssue(url string) bool {
	if url == "" {
		return false
	}
	return strings.Contains(url, "/issues/")
}

func IsPullRequest(url string) bool {
	if url == "" {
		return false
	}
	return strings.Contains(url, "/pull/")
}

func InitClientV4Cursor(sha string) string {
	return fmt.Sprintf("%s 0", sha)
}

func FromClientV4Cursor2Sha(cursor string) string {
	if !strings.Contains(cursor, " ") {
		return cursor
	}
	return strings.Split(cursor, " ")[0]
}
