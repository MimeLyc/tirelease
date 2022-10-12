package model

import (
	"tirelease/internal/entity"
)

type IssuePrRelation struct {
	Major      int
	Minor      int
	Issue      *entity.Issue
	RelatedPrs []entity.PullRequest
}
