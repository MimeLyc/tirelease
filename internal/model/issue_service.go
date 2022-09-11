package model

import (
	"time"
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

func SelectIssuesFixedAfterSprintCheckout(major, minor int, option entity.IssueOption) ([]entity.Issue, error) {
	issuePrRelations, err := GetIssuePrRelations(major, minor, option)
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

func SelectIssuesFixedBeforeSprintCheckout(major, minor int) ([]entity.Issue, error) {
	repos, err := repository.SelectRepo(nil)
	if err != nil {
		return nil, err
	}

	// Get all Issues fixed before sprint checkout.
	masterIssues := make([]entity.Issue, 0)
	for _, repo := range *repos {
		sprintMeta, err := NewSprintMeta(major, minor, repo)
		if err != nil {
			return nil, err
		}

		startTime := *sprintMeta.StartTime
		checkoutTime := time.Now()
		if sprintMeta.CheckoutCommitTime != nil {
			checkoutTime = *sprintMeta.CheckoutCommitTime
		}

		issues, err := repository.SelectIssue(
			&entity.IssueOption{
				State:        "closed",
				CloseTime:    startTime,
				CloseTimeEnd: checkoutTime,
			},
		)
		if err != nil {
			return nil, err
		}
		masterIssues = append(masterIssues, *issues...)

	}
	return masterIssues, nil
}
