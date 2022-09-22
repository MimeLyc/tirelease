package model

import (
	"tirelease/internal/entity"
)

type PullRequest struct {
	*entity.PullRequest
}

func NewPullRequest(entityPr entity.PullRequest) PullRequest {
	return PullRequest{
		PullRequest: &entityPr,
	}
}

func ParseToPullRequest(entityPrs []entity.PullRequest) []PullRequest {
	prs := make([]PullRequest, 0)
	for _, entityPr := range entityPrs {
		prs = append(prs, NewPullRequest(entityPr))
	}
	return prs
}
