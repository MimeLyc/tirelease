package dto

import "tirelease/internal/entity"

type SprintMetaRequest struct {
	ID               *int64
	MinorVersionName *string
	Major            *int `json:"major" form:"major"`
	Minor            *int `json:"minor" form:"minor"`

	entity.ListOption
}
