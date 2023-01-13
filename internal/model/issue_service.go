package model

import (
	"time"
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

	issueBuilder := IssueBuilder{}

	issueBuilder = issueBuilder.Option(
		&entity.IssueOption{
			IssueIDs: issueIds,
		}, nil,
	).Command(
		&TriageBuildCommand{
			NeedTriages: true,
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

		issueBuilder := IssueBuilder{}
		issueBuilder = issueBuilder.Option(&issueOption, nil).Command(
			&TriageBuildCommand{
				NeedTriages: true,
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
