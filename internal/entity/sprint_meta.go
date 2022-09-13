package entity

import (
	"time"

	"gorm.io/gorm"
)

type SprintMeta struct {
	gorm.Model
	MinorVersionName     string
	Major                int
	Minor                int
	RepoID               int
	Repo                 Repo `gorm:"foreignKey:RepoID"`
	StartTime            *time.Time
	CheckoutCommitTime   *time.Time
	BeforeStartCommitSha string
	StartCommitSha       string
	CheckoutCommitSha    string
}

type SprintMetaOption struct {
	ID                 *int64
	MinorVersionName   *string
	Major              *int
	Minor              *int
	Repo               *Repo `gorm:"foreignKey:RepoID"`
	StartTime          *time.Time
	CheckoutCommitTime *time.Time

	ListOption
}
