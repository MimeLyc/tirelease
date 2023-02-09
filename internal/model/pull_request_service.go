package model

import (
	"fmt"
	"time"
	"tirelease/commons/git"
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

// TODO refactor to find from sprint model
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

// TODO refactor to find from sprint model
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
		prs, err := PullRequestCmd{
			PROptions: &entity.PullRequestOption{
				Merged:       &isMerged,
				MergeTime:    &startTime,
				MergeTimeEnd: &checkoutTime,
				Owner:        repo.Owner,
				Repo:         repo.Repo,
			},
			IsDefaultBaseBranch: true,
		}.Build()
		if err != nil {
			return nil, err
		}
		masterPrs = append(masterPrs, prs...)
	}
	return masterPrs, nil
}

// TODO refactor to find from sprint model
func SelectMergePrsAfterSprintCheckout(major, minor int) ([]PullRequest, error) {
	sprintName := ComposeVersionMinorNameByNumber(major, minor)
	branchName := fmt.Sprintf("%s%s", git.ReleaseBranchPrefix, sprintName)
	isMerged := true
	prs, err := PullRequestCmd{
		PROptions: &entity.PullRequestOption{
			Merged:     &isMerged,
			BaseBranch: branchName,
		},
	}.Build()
	if err != nil {
		return nil, err
	}

	return prs, nil
}

func IsPrsAllMerged(prs []PullRequest) bool {
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

func extractPrIds(prs []PullRequest) []string {
	result := make([]string, 0)
	for _, pr := range prs {
		result = append(result, pr.PullRequestID)
	}
	return result
}
