package model

import (
	"tirelease/internal/entity"
)

func ExtractPrIdsByIssueId(relations []entity.IssuePrRelation, issueId string) []string {
	prids := make([]string, 0)
	for _, relation := range relations {
		if relation.IssueID == issueId {
			prids = append(prids, relation.PullRequestID)
		}
	}
	return prids
}

func FilterIssuePrRelationByIssueAndVersion(relations []IssuePrRelation, issueID string, major, minor int) *IssuePrRelation {
	for _, relation := range relations {
		if relation.Issue.IssueID == issueID &&
			relation.Major == major && relation.Minor == minor {
			return &relation
		}
	}
	return nil
}
