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

func extractAssigneeGhLoginsFromIssues(issues *[]entity.Issue) (result []string) {
	for _, issue := range *issues {
		for _, assignee := range *issue.Assignees {
			if login := assignee.GetLogin(); login != "" {
				result = append(result, login)
			}
		}
	}
	return
}

func extractAuthorGhLoginsFromIssues(issues *[]entity.Issue) []string {
	logins := make([]string, 0)
	for _, issue := range *issues {
		logins = append(logins, *&issue.AuthorGHLogin)
	}
	return logins
}

func composeAssignees(issue entity.Issue, loginEmployeeMap map[string]User) []User {
	assignedUsers := make([]User, 0)
	assignees := issue.Assignees
	for _, assignee := range *assignees {
		assignedUsers = append(assignedUsers, loginEmployeeMap[assignee.GetLogin()])
	}

	return assignedUsers
}

func extractIssueIdsFromIssueModels(issues []Issue) []string {
	result := make([]string, 0)

	for _, issue := range issues {
		result = append(result, issue.IssueID)
	}

	return result
}
