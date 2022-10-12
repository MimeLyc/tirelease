package model

import (
	"tirelease/internal/entity"
)

func extractPrIds(relations []entity.IssuePrRelation) []string {
	prids := make([]string, 0)
	for _, relation := range relations {
		prids = append(prids, relation.PullRequestID)
	}
	return prids
}

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
		relation := relation
		if relation.Issue.IssueID == issueID && relation.Major == major && relation.Minor == minor {
			return &relation
		}
	}
	return nil
}

func getPRsByIssueRelation(relations []entity.IssuePrRelation, issueID string, prs *[]entity.PullRequest) []entity.PullRequest {
	result := make([]entity.PullRequest, 0)
	for _, relation := range relations {
		if relation.IssueID != issueID {
			continue
		}

		for _, pr := range *prs {
			if relation.PullRequestID != pr.PullRequestID {
				continue
			}

			result = append(result, pr)
			break
		}
	}
	return result
}
