package model

import (
	"time"
	"tirelease/commons/utils"
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

func SelectIssuesAfterSprintCheckout(major, minor int, issueOption entity.IssueOption) ([]Issue, error) {
	issuePrRelations, err := SelectIssuePrRelationsByVersion(major, minor, issueOption, true)
	if err != nil {
		return nil, err
	}

	issueEntities := make([]entity.Issue, 0)
	for _, relation := range issuePrRelations {
		if relation.Issue.State == "closed" && IsPrsAllMerged(relation.RelatedPrs) {
			relation.Issue.IsFixed = true
		}

		issueEntities = append(issueEntities, *relation.Issue)
	}

	issueIds := extractIssueIdsFromIssues(issueEntities)

	issueBuilder := IssueCmd{}

	issueBuilder = issueBuilder.Option(
		&entity.IssueOption{
			IssueIDs: issueIds,
		}, nil,
	).Command(
		&TriageBuildCommand{
			WithTriages: true,
		},
	)
	issues, err := issueBuilder.BuildArray()
	if err != nil {
		return nil, err
	}

	for _, issue := range issues {
		for _, entity := range issueEntities {
			if entity.IssueID != issue.IssueID {
				continue
			}

			entity := entity
			issue.Issue = entity
		}
	}

	return issues, nil
}

func SelectIssuesBeforeSprintCheckout(major, minor int, issueOption entity.IssueOption) ([]Issue, error) {
	repos, err := repository.SelectRepo(nil)
	if err != nil {
		return nil, err
	}

	// Get all Issues fixed before sprint checkout.
	masterIssues := make([]Issue, 0)
	for _, repo := range *repos {
		sprintMeta, err := NewSprintMeta(major, minor, repo)
		if err != nil {
			// skip error because there are some repos not checking out release branchs
			continue
		}

		startTime := *sprintMeta.StartTime
		checkoutTime := time.Now()
		if sprintMeta.CheckoutCommitTime != nil {
			checkoutTime = *sprintMeta.CheckoutCommitTime
		}

		issueOption.CloseTime = startTime
		issueOption.CloseTimeEnd = checkoutTime
		issueOption.Owner = repo.Owner
		issueOption.Repo = repo.Repo
		issueOption.State = "closed"

		issueBuilder := IssueCmd{}
		issueBuilder = issueBuilder.Option(&issueOption, nil).Command(
			&TriageBuildCommand{
				WithTriages: true,
			},
		)
		issues, err := issueBuilder.BuildArray()
		if err != nil {
			return nil, err
		}
		for _, issue := range issues {
			issue := issue
			issue.IsFixed = true
			masterIssues = append(masterIssues, issue)
		}

	}
	return masterIssues, nil
}

func FilterIssueByID(issues []entity.Issue, issueID string) *entity.Issue {
	for _, issue := range issues {
		if issue.IssueID == issueID {
			return &issue
		}
	}
	return nil
}

func extractIssueIdsFromIssues(issues []entity.Issue) []string {
	issueIds := make([]string, 0)
	for _, issue := range issues {
		issueIds = append(issueIds, issue.IssueID)
	}
	return issueIds

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

func filterIssuesByIssueIds(issues []Issue, issueIds []string) []Issue {
	result := make([]Issue, 0)
	for _, issue := range issues {
		if utils.Contains(issueIds, issue.IssueID) {
			result = append(result, issue)
		}
	}
	return result
}
