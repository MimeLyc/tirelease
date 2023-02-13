package git

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

const (
	org        = "[a-zA-Z0-9][a-zA-Z0-9-]{0,38}"
	repo       = "[a-zA-Z0-9-_]{1,100}"
	orgGroup   = "org"
	repoGroup  = "repo"
	issueGroup = "issue_number"

	issueUrl = "(https|http)://github\\.com/(?P<org>%s)/(?P<repo>%s)/issues/(?P<issue_number>[1-9]\\d*)"
)

var issueUrlRegex = fmt.Sprintf(issueUrl, org, repo)

func ParseIssueUrl(url string) (string, string, int, error) {
	compile, err := regexp.Compile(issueUrlRegex)
	if err != nil {
		return "", "", 0, err
	}

	matches := compile.FindStringSubmatch(url)
	groupNames := compile.SubexpNames()
	if matches == nil {
		return "", "", 0, errors.New("invalid issue url")
	}

	org := ""
	repo := ""
	issue := 0
	for i, match := range matches {
		groupName := groupNames[i]
		if groupName == orgGroup {
			org = match
		} else if groupName == repoGroup {
			repo = match
		} else if groupName == issueGroup {
			issue, err = strconv.Atoi(match)
			if err != nil {
				return "", "", 0, err
			}
		}
	}
	if org == "" || repo == "" || issue == 0 {
		return "", "", 0, errors.New("invalid issue url")
	}

	return org, repo, issue, nil

}
