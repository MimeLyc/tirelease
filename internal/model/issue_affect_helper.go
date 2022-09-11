package model

import "tirelease/internal/entity"

func ExtractIssueIDs(affects []entity.IssueAffect) []string {
	issueIds := make([]string, 0)
	for _, affect := range affects {
		issueIds = append(issueIds, affect.IssueID)
	}
	return issueIds
}
