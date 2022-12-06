package service

import (
	"fmt"
	"strconv"
	"strings"
	"tirelease/commons/git"
	"tirelease/internal/entity"
	"tirelease/internal/model"
	"tirelease/internal/repository"

	"regexp"

	"github.com/google/go-github/v41/github"
)

func WebhookRefreshRef(ref, refType string, repo github.Repository) error {
	if !validateRef(ref, refType) {
		return nil
	}

	versionItems := strings.Split(strings.Split(ref, "-")[1], ".")
	major, _ := strconv.Atoi(versionItems[0])
	minor, _ := strconv.Atoi(versionItems[1])
	return refreshSprint(major, minor, *repo.Owner.Login, *repo.Name)
}

func refreshSprint(major, minor int, owner, repo string) error {
	repoEntity, err := repository.SelectRepo(
		&entity.RepoOption{
			Owner: owner,
			Repo:  repo,
		},
	)
	if err != nil {
		return err
	}
	if len(*repoEntity) == 0 {
		return nil
	}

	targetRepo := &(*repoEntity)[0]
	entitySprint, _ := repository.SelectSprintMetaUnique(
		&entity.SprintMetaOption{
			Major: &major,
			Minor: &minor,
			Repo:  targetRepo,
		},
	)

	if entitySprint != nil {
		return nil
	}

	sprintName := ComposeVersionMinorNameByNumber(major, minor)
	sprint := entity.SprintMeta{
		MinorVersionName: sprintName,
		Major:            major,
		Minor:            minor,
		Repo:             *targetRepo,
	}

	checkoutCommit, err := model.GetCheckoutCommit(owner, repo, sprintName)
	sprint.CheckoutCommitSha = checkoutCommit.Oid
	sprint.CheckoutCommitTime = &checkoutCommit.CommittedTime

	preSprintName, err := model.GetLastMinorVersionName(major, minor)
	preCheckoutCommit, err := model.GetCheckoutCommit(owner, repo, preSprintName)
	sprint.StartCommitSha = preCheckoutCommit.Oid
	sprint.StartTime = &preCheckoutCommit.CommittedTime

	return repository.CreateOrUpdateSprint(&sprint)
}

func validateRef(ref, refType string) bool {
	ltsPattern := fmt.Sprintf("[.*\\/]*%s[0-9]*\\.[0-9]*$", git.ReleaseBranchPrefix)
	ok, err := regexp.MatchString(ltsPattern, ref)
	if err != nil {
		return false
	}
	return refType == git.GitCreateEventRefTypeBranch && ok
}
