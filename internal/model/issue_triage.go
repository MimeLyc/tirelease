package model

import "tirelease/internal/entity"

type IssueTriage struct {
	Issue          entity.Issue
	MasterPrs      []PullRequest
	VersionTriages *[]IssueVersionTriage
}
