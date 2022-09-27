package model

import "tirelease/internal/entity"

func extractIssueIDsFromTriage(triages []entity.VersionTriage) []string {
	issueIDs := make([]string, 0)
	for _, triage := range triages {
		issueIDs = append(issueIDs, triage.IssueID)
	}
	return issueIDs
}
