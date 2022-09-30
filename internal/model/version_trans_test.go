package model

import (
	"fmt"
	"strings"
	"testing"
	"time"
	"tirelease/commons/configs"
	"tirelease/commons/database"
	"tirelease/commons/git"
	"tirelease/internal/entity"
	"tirelease/internal/repository"

	"github.com/stretchr/testify/assert"
)

func TestVersion2Frozen(t *testing.T) {
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	configs.LoadConfig("../../config.yaml")
	config := configs.Config
	database.Connect(config)

	versionName := "1.1.1"

	versionEntity := entity.ReleaseVersion{
		Name:     versionName,
		Major:    1,
		Minor:    1,
		Patch:    1,
		Addition: "",
		Type:     entity.ReleaseVersionTypePatch,
		Status:   entity.ReleaseVersionStatusUpcoming,
	}

	err := repository.CreateReleaseVersion(&versionEntity)
	assert.Nil(t, err)

	version := Parse2ReleaseVersion(versionEntity)
	version.ChangeStatus(entity.ReleaseVersionStatusFrozen)
	assert.Equal(t, entity.ReleaseVersionStatusFrozen, version.Version.Status)
	_, err = repository.DeleteReleaseVersionByName(versionName)
	assert.Nil(t, err)

	// try change version status with mocked data
	err = repository.CreateReleaseVersion(&versionEntity)
	assert.Nil(t, err)

}

func TestVersion2Frozen2(t *testing.T) {
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	configs.LoadConfig("../../config.yaml")
	config := configs.Config
	database.Connect(config)

	versionName := "1.1.1"

	versionEntity := entity.ReleaseVersion{
		Name:     versionName,
		Major:    1,
		Minor:    1,
		Patch:    1,
		Addition: "",
		Type:     entity.ReleaseVersionTypePatch,
		Status:   entity.ReleaseVersionStatusUpcoming,
	}

	// try change version status with mocked data
	err := repository.CreateReleaseVersion(&versionEntity)
	assert.Nil(t, err)
	createMockTriage(versionName)
	version := Parse2ReleaseVersion(versionEntity)
	version.ChangeStatus(entity.ReleaseVersionStatusFrozen)
	assert.Equal(t, entity.ReleaseVersionStatusFrozen, version.Version.Status)
	_, err = repository.DeleteReleaseVersionByName(versionName)
	assert.Nil(t, err)
	deleteMockTriage()
}

func TestVersion2Release(t *testing.T) {
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	configs.LoadConfig("../../config.yaml")
	config := configs.Config
	database.Connect(config)

	versionName := "1.1.1"

	versionEntity := entity.ReleaseVersion{
		Name:     versionName,
		Major:    1,
		Minor:    1,
		Patch:    1,
		Addition: "",
		Type:     entity.ReleaseVersionTypePatch,
		Status:   entity.ReleaseVersionStatusUpcoming,
	}

	// try change version status with mocked data
	err := repository.CreateReleaseVersion(&versionEntity)
	assert.Nil(t, err)
	createMockTriage(versionName)
	version := Parse2ReleaseVersion(versionEntity)
	version.ChangeStatus(entity.ReleaseVersionStatusFrozen)
	assert.Equal(t, entity.ReleaseVersionStatusFrozen, version.Version.Status)
	version.ChangeStatus(entity.ReleaseVersionStatusReleased)
	assert.Equal(t, entity.ReleaseVersionStatusReleased, version.Version.Status)

	_, err = repository.DeleteReleaseVersionByName(versionName)
	assert.Nil(t, err)
	_, err = repository.DeleteReleaseVersionByName("1.1.2")
	assert.Nil(t, err)
	deleteMockTriage()
}

func TestVersion2Cancelled(t *testing.T) {
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	configs.LoadConfig("../../config.yaml")
	config := configs.Config
	database.Connect(config)

	versionName := "1.1.1"

	versionEntity := entity.ReleaseVersion{
		Name:     versionName,
		Major:    1,
		Minor:    1,
		Patch:    1,
		Addition: "",
		Type:     entity.ReleaseVersionTypePatch,
		Status:   entity.ReleaseVersionStatusUpcoming,
	}

	// try change version status with mocked data
	err := repository.CreateReleaseVersion(&versionEntity)
	assert.Nil(t, err)
	createMockTriage(versionName)
	version := Parse2ReleaseVersion(versionEntity)
	version.ChangeStatus(entity.ReleaseVersionStatusReleased)
	assert.Equal(t, entity.ReleaseVersionStatusReleased, version.Version.Status)
	version.ChangeStatus(entity.ReleaseVersionStatusCancelled)
	assert.Equal(t, entity.ReleaseVersionStatusCancelled, version.Version.Status)

	_, err = repository.DeleteReleaseVersionByName(versionName)
	assert.Nil(t, err)
	_, err = repository.DeleteReleaseVersionByName("1.1.2")
	assert.Nil(t, err)
	deleteMockTriage()
}

var issueId = "mockIssue"
var prId = "mockPr"

func createMockTriage(versionName string) entity.VersionTriage {
	minorVersion := fmt.Sprintf("%s.%s", strings.Split(versionName, ".")[0], strings.Split(versionName, ".")[1])
	releaseBranch := git.ReleaseBranchPrefix + minorVersion

	err := repository.CreateOrUpdatePullRequest(
		&entity.PullRequest{
			PullRequestID: prId,
			BaseBranch:    releaseBranch,
			CreateTime:    time.Now(),
			UpdateTime:    time.Now(),
		},
	)
	fmt.Printf("%v", err)
	err = repository.CreateOrUpdateIssue(
		&entity.Issue{
			IssueID:    issueId,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		},
	)
	fmt.Printf("%v", err)
	err = repository.CreateIssuePrRelation(
		&entity.IssuePrRelation{
			PullRequestID: prId,
			IssueID:       issueId,
			CreateTime:    time.Now(),
			UpdateTime:    time.Now(),
		},
	)
	fmt.Printf("%v", err)
	triage := entity.VersionTriage{
		IssueID:      issueId,
		VersionName:  versionName,
		TriageResult: entity.VersionTriageResultAccept,
		CreateTime:   time.Now(),
		UpdateTime:   time.Now(),
	}
	repository.CreateVersionTriage(&triage)
	return triage
}

func deleteMockTriage() {
	_, err := repository.DeletePrsByPRId(prId)
	fmt.Printf("%v", err)
	_, err = repository.DeleteIssueByIssueID(issueId)
	fmt.Printf("%v", err)
	_, err = repository.DeleteRelationByIssueId(issueId)
	fmt.Printf("%v", err)
	repository.DeleteVersionTriagesByIssueId(issueId)
}
