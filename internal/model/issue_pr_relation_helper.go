package model

import "tirelease/internal/entity"

func ExtractPrIdsByIssueId(relations []entity.IssuePrRelation, issueId string) []string {
	prids := make([]string, 0)
	for _, relation := range relations {
		if relation.IssueID == issueId {
			prids = append(prids, relation.PullRequestID)
		}
	}
	return prids
}
