package model

import (
	"tirelease/commons/utils"
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

type PullRequestCmd struct {
	// Filter options
	PROptions *entity.PullRequestOption
	// -- filter by additional pr info
	IsDefaultBaseBranch bool
	// -- filter by related issues
	ByRelatedIssue bool
	IssueIds       []string

	// Filed options
}

func (cmd PullRequestCmd) Build() ([]PullRequest, error) {
	// Find target pr entites
	prs, err := cmd.findPullRequestEntities()
	if err != nil {
		return nil, nil
	}

	// Complete PullRequest fields
	ghLogins := extractAuthorGhLoginsFromPrs(prs)

	userMap, err := UserBuilder{}.BuildUsersByGhLogins(ghLogins)
	if err != nil {
		return nil, err
	}

	result := make([]PullRequest, 0)
	for _, pr := range prs {
		pr := pr
		result = append(result, PullRequest{
			PullRequest: &pr,
			Author:      userMap[pr.AuthorGhLogin],
		})
	}

	return result, nil
}

func (cmd PullRequestCmd) findPullRequestEntities() ([]entity.PullRequest, error) {
	result := make([]entity.PullRequest, 0)
	option := cmd.PROptions
	if option == nil {
		option = new(entity.PullRequestOption)
	}

	if cmd.ByRelatedIssue {
		prIdRange, err := cmd.getIssuesRelatedPrIds()
		if err != nil {
			return nil, err
		}
		if len(prIdRange) == 0 {
			return nil, nil
		}
		if len(option.PullRequestIDs) > 0 {
			option.PullRequestIDs = utils.Intersects(prIdRange, option.PullRequestIDs)
		} else {
			option.PullRequestIDs = prIdRange
		}
	}

	if cmd.IsDefaultBaseBranch {
		option.BaseBranch = "master"
		tmp, err := repository.SelectPullRequest(option)
		if err != nil {
			return nil, err
		}
		result = append(result, *tmp...)

		option.BaseBranch = "main"
		tmp, err = repository.SelectPullRequest(option)
		if err != nil {
			return nil, err
		}
		result = append(result, *tmp...)
	} else {
		tmp, err := repository.SelectPullRequest(option)
		if err != nil {
			return nil, err
		}
		result = *tmp
	}

	return result, nil

}

func extractAuthorGhLoginsFromPrs(prs []entity.PullRequest) []string {
	logins := make([]string, 0)
	for _, pr := range prs {
		logins = append(logins, *&pr.AuthorGhLogin)
	}
	return logins
}

func (cmd PullRequestCmd) getIssuesRelatedPrIds() ([]string, error) {
	prIdRange := make([]string, 0)
	issuePrOption := &entity.IssuePrRelationOption{
		IssueIDs: cmd.IssueIds,
	}
	issuePrRelations, err := repository.SelectIssuePrRelation(issuePrOption)
	if nil != err {
		return nil, err
	}

	if len(*issuePrRelations) > 0 {
		for i := range *issuePrRelations {
			issuePrRelation := (*issuePrRelations)[i]
			prIdRange = append(prIdRange, issuePrRelation.PullRequestID)
		}
	}

	return prIdRange, nil
}
