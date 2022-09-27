package model

import (
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

func getRelatedPullRequests(relatedPrs []entity.IssuePrRelation) ([]entity.PullRequest, error) {
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

func GetIssuePrRelation(issueIDs []string) ([]entity.IssuePrRelation, error) {
	issuePrRelationAll := make([]entity.IssuePrRelation, 0)

	if len(issueIDs) > 0 {
		issuePrRelation := &entity.IssuePrRelationOption{
			IssueIDs: issueIDs,
		}
		issuePrRelationAlls, err := repository.SelectIssuePrRelation(issuePrRelation)
		if nil != err {
			return nil, err
		}
		issuePrRelationAll = append(issuePrRelationAll, (*issuePrRelationAlls)...)
	}

	return issuePrRelationAll, nil
}
