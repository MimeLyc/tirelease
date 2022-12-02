package model

import (
	"fmt"
	"time"
	"tirelease/commons/git"
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

func SelectMergedPrsOfSprint(major, minor int) ([]PullRequest, error) {
	prs, err := SelectMergedPrsBeforeSprintCheckout(major, minor)
	if err != nil {
		return nil, err
	}

	branchPrs, err := SelectMergePrsAfterSprintCheckout(major, minor)
	if err != nil {
		return nil, err
	}
	prs = append(prs, branchPrs...)
	return prs, nil
}

// Select the prs merge into master/main branch before the target sprint is checked out.
func SelectMergedPrsBeforeSprintCheckout(major, minor int) ([]PullRequest, error) {
	repos, err := repository.SelectRepo(nil)
	if err != nil {
		return nil, err
	}

	// Get all merged prs before sprint checkout.
	masterPrs := make([]PullRequest, 0)
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

		isMerged := true
		prs, err := repository.SelectPullRequest(
			&entity.PullRequestOption{
				Merged:       &isMerged,
				MergeTime:    &startTime,
				MergeTimeEnd: &checkoutTime,
				Owner:        repo.Owner,
				Repo:         repo.Repo,
				BaseBranch:   "master",
			},
		)
		if err != nil {
			return nil, err
		}
		masterPrs = append(masterPrs, ParseToPullRequest(*prs)...)
		prs, err = repository.SelectPullRequest(
			&entity.PullRequestOption{
				Merged:       &isMerged,
				MergeTime:    &startTime,
				MergeTimeEnd: &checkoutTime,
				Owner:        repo.Owner,
				Repo:         repo.Repo,
				BaseBranch:   "main",
			},
		)
		if err != nil {
			return nil, err
		}

		masterPrs = append(masterPrs, ParseToPullRequest(*prs)...)
	}
	return masterPrs, nil
}

func SelectMergePrsAfterSprintCheckout(major, minor int) ([]PullRequest, error) {
	sprintName := ComposeVersionMinorNameByNumber(major, minor)
	branchName := fmt.Sprintf("%s%s", git.ReleaseBranchPrefix, sprintName)
	isMerged := true
	entityPrs, err := repository.SelectPullRequest(
		&entity.PullRequestOption{
			Merged:     &isMerged,
			BaseBranch: branchName,
		},
	)
	if err != nil {
		return nil, err
	}

	return ParseToPullRequest(*entityPrs), nil
}

func IsPrsAllMerged(prs []entity.PullRequest) bool {
	mergedCnt := 0
	closedCnt := 0

	for _, pr := range prs {
		if pr.Merged == true {
			mergedCnt++
			continue
		}

		if pr.State == "closed" {
			closedCnt++
		}
	}

	return mergedCnt+closedCnt == len(prs) && mergedCnt > 0
}

func SelectRelatedPrsInMaster(issueID string) ([]entity.PullRequest, error) {
	result, err := SelectRelatedPrs("master", issueID)
	if err != nil {
		return nil, err
	}
	mainPrs, err := SelectRelatedPrs("main", issueID)
	result = append(result, mainPrs...)

	return result, err
}

func SelectRelatedPrs(releaseBranch, issueID string) ([]entity.PullRequest, error) {
	issuePrOption := &entity.IssuePrRelationOption{
		IssueID: issueID,
	}
	issuePrRelations, err := repository.SelectIssuePrRelation(issuePrOption)
	if nil != err {
		return nil, err
	}

	pullRequestIDs := make([]string, 0)
	result := make([]entity.PullRequest, 0)

	if len(*issuePrRelations) > 0 {
		for i := range *issuePrRelations {
			issuePrRelation := (*issuePrRelations)[i]
			pullRequestIDs = append(pullRequestIDs, issuePrRelation.PullRequestID)
		}
		pullRequestOption := &entity.PullRequestOption{
			PullRequestIDs: pullRequestIDs,
			BaseBranch:     releaseBranch,
		}
		pullRequestAlls, err := repository.SelectPullRequest(pullRequestOption)
		if nil != err {
			return nil, err
		}
		result = append(result, (*pullRequestAlls)...)
	}

	return result, nil
}

func extractPrIds(prs []PullRequest) []string {
	result := make([]string, 0)
	for _, pr := range prs {
		result = append(result, pr.PullRequestID)
	}
	return result
}
