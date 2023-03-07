package entity

import (
	"strings"
	"time"
)

type Hotfix struct {
	ID int64 `json:"id"`

	CreateTime        time.Time  `json:"create_time,omitempty"`
	UpdateTime        time.Time  `json:"update_time,omitempty"`
	ActualReleaseTime *time.Time `json:"actual_release_time,omitempty"`

	Name            string `json:"name,omitempty" uri:"name"`
	Customer        string `json:"customer,omitempty"`
	BaseVersionName string `json:"base_version,omitempty"`

	CreatorEmail string       `json:"creator_email,omitempty"`
	Status       HotfixStatus `json:"status,omitempty"`

	Platform     string `json:"platform,omitempty"`
	PassPrecheck bool   `json:"pass_precheck,omitempty" `

	IsDebug    bool `json:"is_debug,omitempty"`
	IsOnHotfix bool `json:"is_on_hotfix,omitempty"`

	HasControlSwitch bool    `json:"has_control_switch,omitempty"`
	RollbackMethod   *string `json:"rollback_method,omitempty"`
	TriggerReason    *string `json:"trigger_reason,omitempty"`

	*Oncall
	*ArtifactConfig

	IsDeleted bool `json:"is_deleted,omitempty"`
}

// DB-Table
func (Hotfix) TableName() string {
	return "hotfix"
}

type Oncall struct {
	OncallPrefix string `json:"oncall_prefix,omitempty"`
	OncallID     string `json:"oncall_id,omitempty"`
	OncallUrl    string `json:"oncall_url,omitempty"`
}

type ArtifactConfig struct {
	ArtifactArchs    string `json:"artifact_archs,omitempty"`
	ArtifactEditions string `json:"artifact_editions,omitempty"`
	ArtifactTypes    string `json:"artifact_types,omitempty"`
}

func (hotfix Hotfix) UnserializeArtifactArchs() []string {
	return strings.Split(hotfix.ArtifactArchs, ",")
}

func (hotfix Hotfix) UnserializeArtifactEditions() []string {
	return strings.Split(hotfix.ArtifactEditions, ",")
}

func (hotfix Hotfix) UnserializeArtifactTypes() []string {
	return strings.Split(hotfix.ArtifactTypes, ",")
}

type HotfixOptions struct {
	ID              int64        `json:"id" form:"id"`
	Name            string       `json:"name,omitempty" form:"name" uri:"name"`
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
