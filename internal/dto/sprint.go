package dto

import (
	"tirelease/internal/entity"
	"tirelease/internal/model"
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
	Major        int            `json:"major,omitempty" form:"major"`
	Minor        int            `json:"minor,omitempty" form:"minor"`
	MasterIssues *[]SprintIssue `json:"master_issues,omitempty"`
	BranchIssues *[]SprintIssue `json:"branch_issues,omitempty"`
}

type SprintRequest struct {
	SprintOption entity.SprintMetaOption
}

type SprintResponse struct {
	Major      int                 `json:"major,omitempty" form:"major"`
	Minor      int                 `json:"minor,omitempty" form:"minor"`
	SprintMeta []entity.SprintMeta `json:"sprint_meta,omitempty"`
}

type SprintIssue struct {
	model.Issue
	IsBlock bool `json:"is_block"`
}
