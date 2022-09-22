package model

import (
	"tirelease/internal/entity"
	"tirelease/internal/repository"
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

func ComposePrIssueRelations(prs []PullRequest) ([]PrIssueRelation, error) {
	if len(prs) == 0 {
		return nil, nil
	}

	prids := extractPrIds(prs)
	issuePrRelations, err := repository.SelectIssuePrRelation(
		&entity.IssuePrRelationOption{
			PullRequestIDs: prids,
		},
	)
	if err != nil {
		return nil, err
	}
	issueIds := ExtractIssueIDsFromRelations(*issuePrRelations)
	issues, err := repository.SelectIssue(
		&entity.IssueOption{
			IssueIDs: issueIds,
		},
	)
	if err != nil {
		return nil, err
	}

	result := make([]PrIssueRelation, 0)
	for _, pr := range prs {
		pr := pr
		relation := PrIssueRelation{
			PullRequest: &pr,
		}
		relatedIssues := make([]entity.Issue, 0)
		for _, issuePrRelation := range *issuePrRelations {
			if issuePrRelation.PullRequestID != pr.PullRequestID {
				continue
			}
			for _, issue := range *issues {
				if issuePrRelation.IssueID == issue.IssueID {
					relatedIssues = append(relatedIssues, issue)
				}
			}
		}
		relation.RelatedIssues = relatedIssues
		result = append(result, relation)
	}

	return result, nil
}

func ExtractIssueIDsFromRelations(relations []entity.IssuePrRelation) []string {
	result := make([]string, 0)
	for _, relation := range relations {
		result = append(result, relation.IssueID)
	}

	return result
}
