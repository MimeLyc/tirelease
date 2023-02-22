package entity

import "time"

type Hotfix struct {
	ID int64 `json:"id"`

	CreateTime        time.Time  `json:"create_time,omitempty"`
	UpdateTime        time.Time  `json:"update_time,omitempty"`
	ActualReleaseTime *time.Time `json:"actual_release_time,omitempty"`

	Name            string `json:"name,omitempty"`
	BaseVersionName string `json:"base_version,omitempty"`

	CreatorEmail string       `json:"creator_email,omitempty"`
	Status       HotfixStatus `json:"status,omitempty"`

	IsDeleted bool `json:"is_deleted,omitempty"`
}

// DB-Table
func (Hotfix) TableName() string {
	return "hotfix"
}

type HotfixOptions struct {
	ID              int64        `json:"id" form:"id"`
	Name            string       `json:"name,omitempty" form:"name"`
	BaseVersionName string       `json:"base_version,omitempty" form:"base_version"`
	CreatorEmail    string       `json:"creator_email,omitempty" form:"creator_email"`
	Status          HotfixStatus `json:"status,omitempty" form:"status"`

	IsDeleted bool `json:"is_deleted,omitempty" form:"is_deleted"`

	ListOption
}

// Enum status
type HotfixStatus string

const (
	HotfixStatusInit            = HotfixStatus("init")
	HotfixStatusPendingApproval = HotfixStatus("pending_approval")
	HotfixStatusDenied          = HotfixStatus("denied")
	HotfixStatusUpcoming        = HotfixStatus("upcoming")
	// HotfixStatusQATesting= HotfixStatus("qa_testing")
	HotfixStatusReleased  = HotfixStatus("released")
	HotfixStatusCancelled = HotfixStatus("cancelled")
)
