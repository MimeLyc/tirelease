package model

import (
	"tirelease/commons/git"
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

type IssuePrRelation struct {
	Major      int
	Minor      int
	Issue      *entity.Issue
	RelatedPrs []entity.PullRequest
}

type PrIssueRelation struct {
	*PullRequest
	RelatedIssues []Issue
}

func GetIssuePrRelations(major, minor int, option entity.IssueOption) ([]IssuePrRelation, error) {
	versionName := ComposeVersionMinorNameByNumber(major, minor)
	branchName := git.ReleaseBranchPrefix + versionName
	affects, err := repository.SelectIssueAffect(
		&entity.IssueAffectOption{
			AffectVersion: versionName,
			AffectResult:  entity.AffectResultResultYes,
		},
	)
	if err != nil || len(*affects) == 0 {
		return nil, err
	}

	issueIds := ExtractIssueIDs(*affects)
	option.IssueIDs = issueIds
	issues, err := repository.SelectIssue(&option)
	if err != nil {
		return nil, err
	}
	issuePrRelations, err := repository.SelectIssuePrRelation(
		&entity.IssuePrRelationOption{
			IssueIDs: issueIds,
		},
	)
	if err != nil {
		return nil, err
	}

	result := make([]IssuePrRelation, 0)

	for _, issue := range *issues {
		issue := issue
		prids := ExtractPrIdsByIssueId(*issuePrRelations, issue.IssueID)
		if len(prids) == 0 {
			continue
		}
		prs, err := repository.SelectPullRequest(
			&entity.PullRequestOption{
				BaseBranch:     branchName,
				PullRequestIDs: prids,
			},
		)
		if err != nil {
			return nil, err
		}

		result = append(result, IssuePrRelation{
			Major:      major,
			Minor:      minor,
			Issue:      &issue,
			RelatedPrs: *prs,
		})
	}

	return result, nil
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
	issues, err := IssueBuilder{}.BuildIssues(
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
