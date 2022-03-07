package dto

import (
	"tirelease/internal/entity"
)

// VersionTriage Query Struct
type VersionTriageInfoQuery struct {
	entity.VersionTriageOption
}

// VersionTriage ReturnBack Struct
type VersionTriageInfo struct {
	ReleaseVersion *entity.ReleaseVersion `json:"release_version,omitempty"`
	IsFrozen       bool                   `json:"is_frozen,omitempty"`
	IsAccept       bool                   `json:"is_accept,omitempty"`

	VersionTriage     *entity.VersionTriage `json:"version_triage,omitempty"`
	IssueRelationInfo *IssueRelationInfo    `json:"issue_relation_info,omitempty"`
}

type VersionTriageInfoWrap struct {
	ReleaseVersion     *entity.ReleaseVersion `json:"release_version,omitempty"`
	VersionTriageInfos *[]VersionTriageInfo   `json:"version_triage_infos,omitempty"`
}
