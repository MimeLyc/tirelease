package model

import (
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

type Issue struct {
	entity.Issue
	Assignees []User `json:"assigned_employees,omitempty"`
	Author    User   `json:"author,omitempty"`
}

type IssueBuilder struct {
}

func (builder IssueBuilder) BuildIssues(option *entity.IssueOption) ([]Issue, error) {
	issues, err := repository.SelectIssue(
		option,
	)

	if err != nil {
		return nil, err
	}

	ghLogins := extractAuthorGhLoginsFromIssues(issues)
	ghLogins = append(ghLogins, extractAssigneeGhLoginsFromIssues(issues)...)
	userMap, err := UserBuilder{}.BuildUsersByGhLogins(ghLogins)
	if err != nil {
		return nil, err
	}

	result := make([]Issue, 0)
	for _, issue := range *issues {
		result = append(result, Issue{
			Issue:     issue,
			Assignees: composeAssignees(issue, userMap),
			Author:    userMap[issue.AuthorGHLogin],
		})

	}
	return result, nil
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
