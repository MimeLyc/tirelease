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

func SelectIssuePrRelations(major, minor int, option entity.IssueOption, limitAffect bool) ([]IssuePrRelation, error) {
	versionName := ComposeVersionMinorNameByNumber(major, minor)
	branchName := git.ReleaseBranchPrefix + versionName
	if limitAffect {
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
	}

	issues, err := repository.SelectIssue(&option)
	if err != nil {
		return nil, err
	}

	issueIDs := extractIssueIdsFromIssues(*issues)
	issuePrRelations, err := repository.SelectIssuePrRelation(
		&entity.IssuePrRelationOption{
			IssueIDs: issueIDs,
		},
	)
	if err != nil {
		return nil, err
	}

	result := make([]IssuePrRelation, 0)

	for _, issue := range *issues {
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
