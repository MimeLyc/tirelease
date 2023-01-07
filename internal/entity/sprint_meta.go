package entity

import (
	"time"

	"gorm.io/gorm"
)

type SprintMeta struct {
	gorm.Model
	MinorVersionName   string
	Major              int
	Minor              int
	RepoFullName       string
	Repo               Repo `gorm:"foreignKey:RepoFullName;references:FullName"`
	StartTime          *time.Time
	CheckoutCommitTime *time.Time
	StartCommitSha     string
	CheckoutCommitSha  string
}

type SprintMetaOption struct {
	ID                 *int64     `json:"id,omitempty"`
	MinorVersionName   *string    `json:"minor_version_name,omitempty"`
	Major              *int       `json:"major,omitempty" form:"major"`
	Minor              *int       `json:"minor,omitempty" form:"minor"`
	Repo               *Repo      `gorm:"foreignKey:RepoFullName; references:FullName" json:"repo,omitempty"`
	StartTime          *time.Time `json:"start_time,omitempty"`
	CheckoutCommitTime *time.Time `json:"checkout_commit_time,omitempty"`

	ListOption `json:"list_option,omitempty"`
}

// DB-Table
func (SprintMeta) TableName() string {
	return "sprint_info"
}
