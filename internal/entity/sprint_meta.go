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
	ID                 *int64
	MinorVersionName   *string
	Major              *int
	Minor              *int
	Repo               *Repo `gorm:"foreignKey:RepoFullName; references:FullName"`
	StartTime          *time.Time
	CheckoutCommitTime *time.Time

	ListOption
}

// DB-Table
func (SprintMeta) TableName() string {
	return "sprint_info"
}
