package dto

import (
	"tirelease/internal/model"
)

type HotfixSaveRequest struct {
	model.Hotfix
	ReleaseInfos  []HotfixReleaseInfoRequest `json:"release_infos"`
	OperatorEmail string                     `json:"operator_email"`
}

// Validate method validate the request.
func (r *HotfixSaveRequest) Validate() error {
	return nil
}

type HotfixReleaseInfoRequest struct {
	model.HotfixReleaseInfo
	Repo      string        `json:"repo,omitempty"`
	Issues    []HotfixIssue `json:"issues,omitempty"`
	MasterPrs []HotfixPr    `json:"master_prs,omitempty"`
	BranchPrs []HotfixPr    `json:"branch_prs,omitempty"`
}

type HotfixIssue struct {
	Owner   string `json:"owner,omitempty"`
	Repo    string `json:"repo,omitempty"`
	Number  int    `json:"number,omitempty"`
	HTMLURL string `json:"html_url,omitempty"`
}

type HotfixPr struct {
	Owner   string `json:"owner,omitempty"`
	Repo    string `json:"repo,omitempty"`
	Number  int    `json:"number,omitempty"`
	HTMLURL string `json:"html_url,omitempty"`
}
