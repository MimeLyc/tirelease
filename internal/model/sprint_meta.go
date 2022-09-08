package model

import (
	"time"
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

type SprintMeta struct {
	entity.SprintMeta
}

func (sprint SprintMeta) GetMajorVersion() int {
	major, _, _, _ := ComposeVersionAtom(sprint.MinorVersionName)
	return major
}

func (sprint SprintMeta) GetMinorVersion() int {
	_, minor, _, _ := ComposeVersionAtom(sprint.MinorVersionName)
	return minor
}

// StartTime of a sprint if the checkout time of last sprint.
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

	return GetCheckoutTimeOfSprint(sprint.Repo.Owner, sprint.Repo.Repo, lastMinorVersionName)
}

func (sprint SprintMeta) GetCheckoutTime() (*time.Time, error) {
	major := sprint.Major
	minor := sprint.Minor
	sprintEntity, _ := repository.SelectSprintMetaUnique(
		&entity.SprintMetaOption{
			Major: &major,
			Minor: &minor,
			Repo:  &sprint.Repo,
		},
	)

	// If there is data of last Sprint, just return the stored value
	if sprintEntity != nil && sprintEntity.CheckoutCommitTime != nil {
		return sprintEntity.CheckoutCommitTime, nil
	} else {
		return GetCheckoutTimeOfSprint(sprint.Repo.Owner, sprint.Repo.Repo, ComposeVersionMinorNameByNumber(major, minor))
	}
}
