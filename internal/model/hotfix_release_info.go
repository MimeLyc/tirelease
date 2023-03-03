package model

import "tirelease/internal/entity"

type HotfixReleaseInfo struct {
	entity.HotfixReleaseInfo
	Issues    []Issue
	MasterPrs []PullRequest
}
