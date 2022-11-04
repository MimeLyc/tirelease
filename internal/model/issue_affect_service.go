package model

import "tirelease/internal/entity"

func ExtractIssueIDs(affects []entity.IssueAffect) []string {
	issueIds := make([]string, 0)
	for _, affect := range affects {
		issueIds = append(issueIds, affect.IssueID)
	}
	return issueIds
}

func FilterAffectByIssueIDandMinorVersion(affects []entity.IssueAffect,
	issueId, minorVersion string) *entity.IssueAffect {
	for _, affect := range affects {
		if affect.IssueID == issueId && affect.AffectVersion == minorVersion {
			return &affect
		}
	}
	return nil
}
