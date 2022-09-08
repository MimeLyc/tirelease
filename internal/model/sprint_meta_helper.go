package model

import (
	"fmt"
	"time"
	"tirelease/commons/git"
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

func (sprint SprintMeta) GetStartTime() (*time.Time, error) {
	major := sprint.Major
	minor := sprint.Minor
	lastSprint, err := SelectLastSprint(major, minor, sprint.Repo)
	lastMinorVersionName := ""

	// If there is data of last Sprint, just return the stored value
	if lastSprint != nil && lastSprint.CheckoutCommitTime != nil {
		return lastSprint.CheckoutCommitTime, nil
	} else if minor > 0 {
		lastMinorVersionName = ComposeVersionMinorNameByNumber(major, minor-1)
	} else {
		return nil, err
	}

	// If there is tag named vx.x.0, then trace back the checkout commit from it
	// because it's the nearest commit of checkout commit
	firstTagOfLastSprint := fmt.Sprintf("v%s.0", lastMinorVersionName)
	checkoutCommit, err := GetCheckoutCommitOfRef(sprint.Repo.Owner, sprint.Repo.Repo, firstTagOfLastSprint, git.RefTypeTag)
	if err != nil {
		return nil, err
	}
	if checkoutCommit != nil {
		return &checkoutCommit.CommittedTime, nil
	}

	// If there is **not** tag named vx.x.0, trace back the checkout commit from version branch "release-x.x"
	// because there is no tag on the branch
	branchOfLastSprint := fmt.Sprintf("%s%s", git.ReleaseBranchPrefix, lastMinorVersionName)
	checkoutCommit, err = GetCheckoutCommitOfRef(sprint.Repo.Owner, sprint.Repo.Repo, branchOfLastSprint, git.RefTypeBranch)
	if err != nil {
		return nil, err
	}
	if checkoutCommit != nil {
		return &checkoutCommit.CommittedTime, nil
	}

	// If the checkout commit is not found and the err is nil, the branch of last sprint is not checkouted.
	return nil, fmt.Errorf("Last branch of sprint %s has not been checkouted.", lastMinorVersionName)
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
