package model

import "tirelease/internal/entity"

type IssueTriage struct {
	Issue         entity.Issue
	MasterPrs     []entity.PullRequest
	versionTriage []*IssueVersionTriage
}
