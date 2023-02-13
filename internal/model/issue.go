package model

import (
	"tirelease/internal/entity"
)

type Issue struct {
	entity.Issue
	Assignees      []User                 `json:"assigned_employees,omitempty"`
	Author         User                   `json:"author,omitempty"`
	VersionTriages []entity.VersionTriage `json:"version_triages,omitempty"`
}

func (i Issue) ForcePickTriage(version string, triageResult entity.VersionTriageResult) error {
	triage, err := SelectActiveIssueVersionTriage(version, i.IssueID)
	if err != nil {
		return err
	}

	err = triage.ForceTriagePickStatus(triageResult)
	if err != nil {
		return err
	}
	return CreateOrUpdateVersionTriageInfo(triage, entity.VersionTriageUpdatedVarTriageResult)
}

func (i Issue) PickTriage(version string, triageResult entity.VersionTriageResult) error {
	triage, err := SelectActiveIssueVersionTriage(version, i.IssueID)
	if err != nil {
		return err
	}

	err = triage.TriagePickStatus(triageResult)
	if err != nil {
		return err
	}
	return CreateOrUpdateVersionTriageInfo(triage, entity.VersionTriageUpdatedVarTriageResult)
}
