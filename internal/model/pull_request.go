package model

import (
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

type PullRequest struct {
	*entity.PullRequest
	Author User `json:"author,omitempty"`
}

func NewPullRequest(entityPr entity.PullRequest) PullRequest {
	return PullRequest{
		PullRequest: &entityPr,
	}
}

type PullRequestBuilder struct {
}

func (builder PullRequestBuilder) Build(option *entity.PullRequestOption) ([]PullRequest, error) {
	prs, err := repository.SelectPullRequest(option)
	if err != nil {
		return nil, nil
	}

	ghLogins := extractAuthorGhLoginsFromPrs(prs)

	userMap, err := UserBuilder{}.BuildUsersByGhLogins(ghLogins)
	if err != nil {
		return nil, err
	}

	result := make([]PullRequest, 0)
	for _, pr := range *prs {
		pr := pr
		result = append(result, PullRequest{
			PullRequest: &pr,
			Author:      userMap[pr.AuthorGhLogin],
		})
	}

	return result, nil
}

func extractAuthorGhLoginsFromPrs(prs *[]entity.PullRequest) []string {
	logins := make([]string, 0)
	for _, pr := range *prs {
		logins = append(logins, *&pr.AuthorGhLogin)
	}
	return logins
}
