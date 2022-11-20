package model

import (
	"time"
	"tirelease/commons/git"
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

<<<<<<< HEAD
func SelectIssuesFixedAfterSprintCheckout(major, minor int, option entity.IssueOption) ([]entity.Issue, error) {
	issuePrRelations, err := SelectIssuePrRelationsByVersion(major, minor, option, true)
=======
func SelectBugsAfterSprintCheckout(major, minor int) ([]entity.Issue, error) {
	issueOption := entity.IssueOption{
		TypeLabel: git.BugTypeLabel,
	}
	issuePrRelations, err := SelectIssuePrRelations(major, minor, issueOption, true)
>>>>>>> main
	if err != nil {
		return nil, err
	}

	result := make([]entity.Issue, 0)
	for _, relation := range issuePrRelations {
		if relation.Issue.State == "closed" && IsPrsAllMerged(relation.RelatedPrs) {
			relation.Issue.IsFixed = true
		}

		result = append(result, *relation.Issue)
	}

	return result, nil
}

func SelectFixedBugsBeforeSprintCheckout(major, minor int) ([]entity.Issue, error) {
	repos, err := repository.SelectRepo(nil)
	if err != nil {
		return nil, err
	}

	// Get all Issues fixed before sprint checkout.
	masterIssues := make([]entity.Issue, 0)
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

		issueOption := entity.IssueOption{
			CloseTime:    startTime,
			CloseTimeEnd: checkoutTime,
			State:        "closed",
			TypeLabel:    git.BugTypeLabel,
			Owner:        repo.Owner,
			Repo:         repo.Repo,
		}
		issues, err := repository.SelectIssue(&issueOption)
		if err != nil {
			return nil, err
		}
		for _, issue := range *issues {
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
