package model

import (
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

func SelectIssueTriage(issueId string) (*IssueTriage, error) {
	issue, err := repository.SelectIssueUnique(&entity.IssueOption{
		IssueID: issueId,
	})

	if err != nil {
		return nil, err
	}

	masterPrs, err := SelectRelatedPrsInMaster(issueId)
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
