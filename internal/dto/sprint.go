package dto

import (
	"tirelease/internal/entity"
)

type SprintMetaRequest struct {
	ID               *int64
	MinorVersionName *string
	Major            *int `json:"major" form:"major"`
	Minor            *int `json:"minor" form:"minor"`

	entity.ListOption
}

type SprintIssueRequest struct {
	ID    *int64
	Major *int `json:"major" form:"major"`
	Minor *int `json:"minor" form:"minor"`

	entity.IssueOption
	entity.ListOption
}

type SprintIssuesResponse struct {
	Major        int             `json:"major" form:"major"`
	Minor        int             `json:"minor" form:"minor"`
	MasterIssues *[]entity.Issue `json:"master_issues"`
	BranchIssues *[]entity.Issue `json:"branch_issues"`
}
