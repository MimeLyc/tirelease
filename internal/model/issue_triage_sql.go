package model

import (
	"tirelease/internal/entity"
	"tirelease/internal/store"
)

func SelectIssueTriage(issueId string) (*IssueTriage, error) {
	issue, err := store.SelectIssueUnique(&entity.IssueOption{
		IssueID: issueId,
	})

	if err != nil {
		return nil, err
	}

	masterPrs, err := PullRequestCmd{
		IsDefaultBaseBranch: true,
		ByRelatedIssue:      true,
		IssueIds:            []string{issueId},
	}.Build()

	if err != nil {
		return nil, err
	}

	triages, err := SelectAllTriagesByIssue(*issue)

	if err != nil {
		return nil, err
	}

	return &IssueTriage{
		Issue:          *issue,
		MasterPrs:      masterPrs,
		VersionTriages: &triages,
	}, nil
}
