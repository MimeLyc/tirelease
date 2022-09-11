package model

import (
	"fmt"
	"time"
	"tirelease/commons/git"
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

// StartTime of a sprint if the checkout time of last sprint.
func CalculateStartTimeOfSprint(major, minor int, repo entity.Repo) (*time.Time, error) {
	lastSprint, err := SelectLastSprint(major, minor, repo)
	lastMinorVersionName := ""

	// If there is data of last Sprint, just return the stored value
	if lastSprint != nil && lastSprint.CheckoutCommitTime != nil {
		return lastSprint.CheckoutCommitTime, nil
	} else if minor > 0 {
		lastMinorVersionName = ComposeVersionMinorNameByNumber(major, minor-1)
	} else {
		return nil, err
	}

	return GetCheckoutTimeOfSprint(repo.Owner, repo.Repo, lastMinorVersionName)
}

// StartTime of a sprint if the checkout time of last sprint.
func CalculateCheckoutTimeOfSprint(major, minor int, repo entity.Repo) (*time.Time, error) {
	sprintName := ComposeVersionMinorNameByNumber(major, minor)
	return GetCheckoutTimeOfSprint(repo.Owner, repo.Repo, sprintName)
}

// Format of sprint name: x.x
func GetCheckoutTimeOfSprint(owner, repo, sprintName string) (*time.Time, error) {
	commit, err := GetCheckoutCommit(owner, repo, sprintName)
	if commit == nil {
		return nil, err
	}
	return &commit.CommittedTime, nil
}

func GetCheckoutCommit(owner, repo, sprintName string) (*GitCommit, error) {
	// If there is tag named vx.x.0, then trace back the checkout commit from it
	// because it's the nearest commit of checkout commit
	firstTagOfLastSprint := fmt.Sprintf("v%s.0", sprintName)
	checkoutCommit, err := GetCheckoutCommitOfRef(owner, repo, firstTagOfLastSprint, git.RefTypeTag)
	if err != nil {
		return nil, err
	}
	if checkoutCommit != nil {
		return checkoutCommit, nil
	}

	// If there is **not** tag named vx.x.0, trace back the checkout commit from version branch "release-x.x"
	// because there is no tag on the branch
	branchOfLastSprint := fmt.Sprintf("%s%s", git.ReleaseBranchPrefix, sprintName)
	checkoutCommit, err = GetCheckoutCommitOfRef(owner, repo, branchOfLastSprint, git.RefTypeBranch)
	if err != nil {
		return nil, err
	}
	if checkoutCommit != nil {
		return checkoutCommit, nil
	}

	// If the checkout commit is not found and the err is nil, the branch of last sprint is not checkouted.
	return nil, fmt.Errorf("Last branch of sprint %s has not been checkouted.", sprintName)

}

func SelectLastSprint(major, minor int, repo entity.Repo) (*entity.SprintMeta, error) {
	if minor > 0 {
		lastMinor := minor - 1
		return repository.SelectSprintMetaUnique(
			&entity.SprintMetaOption{
				Major: &major,
				Minor: &lastMinor,
				Repo:  &repo,
			},
		)
	} else {
		lastMajor := major - 1
		return repository.SelectSprintMetaUnique(
			&entity.SprintMetaOption{
				Major: &lastMajor,
				Repo:  &repo,
			},
		)
	}
}
