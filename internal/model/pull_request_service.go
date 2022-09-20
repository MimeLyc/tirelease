package model

import (
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

func IsPrsAllMerged(prs []entity.PullRequest) bool {
	mergedCnt := 0
	closedCnt := 0

	for _, pr := range prs {
		if pr.Merged == true {
			mergedCnt++
			continue
		}

		if pr.State == "closed" {
			closedCnt++
		}
	}

	return mergedCnt+closedCnt == len(prs) && mergedCnt > 0
}

func GetRelatedPrs(releaseBranch, issueID string) ([]entity.PullRequest, error) {
	issuePrOption := &entity.IssuePrRelationOption{
		IssueID: issueID,
	}
	issuePrRelations, err := repository.SelectIssuePrRelation(issuePrOption)
	if nil != err {
		return nil, err
	}

	pullRequestIDs := make([]string, 0)
	result := make([]entity.PullRequest, 0)

	if len(*issuePrRelations) > 0 {
		for i := range *issuePrRelations {
			issuePrRelation := (*issuePrRelations)[i]
			pullRequestIDs = append(pullRequestIDs, issuePrRelation.PullRequestID)
		}
		pullRequestOption := &entity.PullRequestOption{
			PullRequestIDs: pullRequestIDs,
			BaseBranch:     releaseBranch,
		}
		pullRequestAlls, err := repository.SelectPullRequest(pullRequestOption)
		if nil != err {
			return nil, err
		}
		result = append(result, (*pullRequestAlls)...)
	}

	return result, nil
}
