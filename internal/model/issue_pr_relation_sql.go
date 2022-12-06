package model

import (
	"tirelease/commons/git"
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

func SelectIssuePrRelationsByVersion(major, minor int, option entity.IssueOption, limitAffect bool) ([]IssuePrRelation, error) {
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

	prids := extractPrIdsFromIssuePrRelation(*issuePrRelations)
	prs, err := repository.SelectPullRequest(
		&entity.PullRequestOption{
			BaseBranch:     branchName,
			PullRequestIDs: prids,
		},
	)

	result := make([]IssuePrRelation, 0)
	if len(*prs) == 0 {
		return nil, nil
	}

	for _, issue := range *issues {
		issue := issue
		issuePrs := getPRsByIssueRelation(*issuePrRelations, issue.IssueID, prs)
		if len(issuePrs) == 0 {
			continue
		}
		if err != nil {
			return nil, err
		}

		result = append(result, IssuePrRelation{
			Major:      major,
			Minor:      minor,
			Issue:      &issue,
			RelatedPrs: issuePrs,
		})
	}

	return result, nil
}

func selectRelatedPullRequests(relatedPrs []entity.IssuePrRelation) ([]entity.PullRequest, error) {
	pullRequestIDs := make([]string, 0)
	pullRequestAll := make([]entity.PullRequest, 0)

	if len(relatedPrs) == 0 {
		return pullRequestAll, nil
	}

	for i := range relatedPrs {
		issuePrRelation := relatedPrs[i]
		pullRequestIDs = append(pullRequestIDs, issuePrRelation.PullRequestID)
	}
	pullRequestOption := &entity.PullRequestOption{
		PullRequestIDs: pullRequestIDs,
	}
	pullRequestAlls, err := repository.SelectPullRequest(pullRequestOption)
	if nil != err {
		return nil, err
	}
	pullRequestAll = append(pullRequestAll, (*pullRequestAlls)...)

	return pullRequestAll, nil
}

func SelectIssuePrRelationByIds(issueIDs []string) ([]entity.IssuePrRelation, error) {
	issuePrRelationAll := make([]entity.IssuePrRelation, 0)

	if len(issueIDs) > 0 {
		issuePrRelationAlls, err := repository.SelectIssuePrRelation(
			&entity.IssuePrRelationOption{
				IssueIDs: issueIDs,
			},
		)

		if nil != err {
			return nil, err
		}
		issuePrRelationAll = append(issuePrRelationAll, (*issuePrRelationAlls)...)
	}

	return issuePrRelationAll, nil
}
