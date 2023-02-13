package model

import (
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

type IssuePrRelation struct {
	Major      int
	Minor      int
	Issue      *entity.Issue
	RelatedPrs []PullRequest
}

type PrIssueRelation struct {
	*PullRequest
	RelatedIssues []Issue
}

type PrIssueRelationBuilder struct {
	Issues      []Issue
	PullRequest []PullRequest
}

func (builder PrIssueRelationBuilder) BuildByPrs(prs []PullRequest) ([]PrIssueRelation, error) {
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
	issues, err := IssueCmd{}.Option(&entity.IssueOption{
		IssueIDs: issueIds,
	}, nil).BuildArray()

	if err != nil {
		return nil, err
	}

	result := make([]PrIssueRelation, 0)
	for _, pr := range prs {
		pr := pr
		relation := PrIssueRelation{
			PullRequest: &pr,
		}
		relatedIssues := make([]Issue, 0)
		for _, issuePrRelation := range *issuePrRelations {
			if issuePrRelation.PullRequestID != pr.PullRequestID {
				continue
			}
			for _, issue := range issues {
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
