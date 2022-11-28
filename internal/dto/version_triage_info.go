package dto

import (
	"tirelease/internal/entity"
)

// VersionTriage Query Struct
type VersionTriageInfoQuery struct {
	entity.VersionTriageOption

	Version string `json:"version,omitempty" form:"version" uri:"version"`
}

// VersionTriage ReturnBack Struct
type VersionTriageInfo struct {
	ReleaseVersion *entity.ReleaseVersion `json:"release_version,omitempty"`
	IsFrozen       bool                   `json:"is_frozen,omitempty"`
	IsAccept       bool                   `json:"is_accept,omitempty"`

	VersionTriage            *entity.VersionTriage           `json:"version_triage,omitempty"`
	VersionTriageMergeStatus entity.VersionTriageMergeStatus `json:"version_triage_merge_status,omitempty"`

	IssueRelationInfo *IssueRelationInfo `json:"issue_relation_info,omitempty"`
}

type VersionTriageInfoWrap struct {
	ReleaseVersion     *entity.ReleaseVersion `json:"release_version,omitempty"`
	VersionTriageInfos *[]VersionTriageInfo   `json:"version_triage_infos,omitempty"`
}

type VersionTriage struct {
	ReleaseVersion *entity.ReleaseVersion `json:"release_version,omitempty"`
	VersionPrs     *[]entity.PullRequest  `json:"pull_requests"`

	PickTriageResult entity.VersionTriageResult `json:"triage_result"`

	BlockTriageResult entity.BlockVersionReleaseResult `json:"block_version_release"`
	IsBlock           bool                             `json:"is_block"`

	AffectResult entity.AffectResultResult `json:"affect_result,omitempty"`
	IsAffect     bool                      `json:"is_affect,omitempty"`

	Comment     string                          `json:"comment"`
	ChangedItem entity.VersionTriageChangedItem `json:"changed_item"`

	MergeStatus entity.VersionTriageMergeStatus `json:"merge_status,omitempty"`
}
