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

func GetIssuePrRelations(major, minor int, option entity.IssueOption) ([]IssuePrRelation, error) {
	versionName := ComposeVersionMinorNameByNumber(major, minor)
	branchName := git.ReleaseBranchPrefix + versionName
	affects, err := repository.SelectIssueAffect(
		&entity.IssueAffectOption{
			AffectVersion: versionName,
			AffectResult:  entity.AffectResultResultYes,
		},
	)
	if err != nil {
		return nil, err
	}

	issueIds := ExtractIssueIDs(*affects)
	option.IssueIDs = issueIds
	issues, err := repository.SelectIssue(&option)
	issuePrRelations, err := repository.SelectIssuePrRelation(
		&entity.IssuePrRelationOption{
			IssueIDs: issueIds,
		},
	)

	result := make([]IssuePrRelation, 0)

	for _, issue := range *issues {
		prids := ExtractPrIdsByIssueId(*issuePrRelations, issue.IssueID)
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
