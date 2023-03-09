package model

import (
	"fmt"
	"tirelease/internal/entity"
	"tirelease/internal/store"
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

// TODO: fill the start commit sha and checkout sha if needed.
func NewSprintMeta(major, minor int, repo entity.Repo) (SprintMeta, error) {
	sprint := SprintMeta{
		entity.SprintMeta{
			Major: major,
			Minor: minor,
			Repo:  repo,
		},
	}

	entitySprint, _ := store.SelectSprintMetaUnique(
		&entity.SprintMetaOption{
			Major: &major,
			Minor: &minor,
			Repo:  &repo,
		},
	)

	if entitySprint != nil {
		sprint.SprintMeta = *entitySprint
	}

	if sprint.StartTime == nil {
		startTime, err := CalculateStartTimeOfSprint(major, minor, repo)
		if err != nil {
			return sprint, err
		}
		sprint.StartTime = startTime
	}

	if sprint.CheckoutCommitTime == nil {
		checkoutTime, _ := CalculateCheckoutTimeOfSprint(major, minor, repo)
		sprint.CheckoutCommitTime = checkoutTime

	}

	return sprint, nil
}

type BranchNotFoundError struct {
	Branch string
}

func (error *BranchNotFoundError) Error() string {
	return fmt.Sprintf("Branch of sprint %s not found", error.Branch)
}
