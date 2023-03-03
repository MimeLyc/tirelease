package dto

import (
	"fmt"
	"time"
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

func (r *HotfixSaveRequest) Map2Model() (model.Hotfix, error) {
	hotfix := r.Hotfix
	if hotfix.Name == "" && hotfix.BaseVersionName != "" && hotfix.OncallID != "" {
		today := time.Now().Format("20060102")
		hotfix.Name = fmt.Sprintf(
			"%s-%s-%s-%s",
			today,
			hotfix.BaseVersionName,
			hotfix.OncallPrefix,
			hotfix.OncallID,
		)
	}

	// TODO fill hotfix release info(issue, pr) from repository.
	return hotfix, nil
}

type HotfixReleaseInfoRequest struct {
	Repo                string        `json:"repo,omitempty"`
	BasedReleaseVersion string        `json:"based_release_version,omitempty"`
	BasedCommitSHA      string        `json:"based_commit_sha,omitempty"`
	Issues              []HotfixIssue `json:"issues,omitempty"`
	MasterPrs           []HotfixPr    `json:"master_prs,omitempty"`
}

type HotfixIssue struct {
	Org    string `json:"org,omitempty"`
	Repo   string `json:"repo,omitempty"`
	Number int    `json:"number,omitempty"`
	Url    string `json:"url,omitempty"`
}

type HotfixPr struct {
	Org    string `json:"org,omitempty"`
	Repo   string `json:"repo,omitempty"`
	Number int    `json:"number,omitempty"`
	Url    string `json:"url,omitempty"`
}
