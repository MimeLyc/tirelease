package model

import (
	"time"
	"tirelease/commons/git"
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

func SelectIssuesFixedAfterSprintCheckout(major, minor int, option entity.IssueOption) ([]entity.Issue, error) {
	issuePrRelations, err := SelectIssuePrRelations(major, minor, option, true)
	if err != nil {
		return nil, err
	}

	result := make([]entity.Issue, 0)
	for _, relation := range issuePrRelations {
		if IsPrsAllMerged(relation.RelatedPrs) {
			result = append(result, *relation.Issue)
		}
	}

	return result, nil
}

func SelectBugsFixedBeforeSprintCheckout(major, minor int) ([]entity.Issue, error) {
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

		issues, err := repository.SelectIssue(
			&entity.IssueOption{
				State:        "closed",
				TypeLabel:    git.BugTypeLabel,
				CloseTime:    startTime,
				CloseTimeEnd: checkoutTime,
				Owner:        repo.Owner,
				Repo:         repo.Repo,
			},
		)
		if err != nil {
			return nil, err
		}
		masterIssues = append(masterIssues, *issues...)

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
