package service

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"testing"
	"tirelease/commons/feishu"
	"tirelease/commons/git"
	"tirelease/internal/entity"
	"tirelease/internal/store"
	"tirelease/utils/configs"

	"github.com/stretchr/testify/assert"
)

func TestNotifySprintIssueMetrics(t *testing.T) {
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	config := configs.NewConfig(
		"../../"+configs.TestConfig,
		"../../"+configs.TestSecretConfig)
	store.NewStore(config)
	// FeishuAPP
	feishu.SetFeishuApp(config.FeiShu.AppId, config.FeiShu.AppSecret)

	major := 6
	minor := 3
	err := NotifySprintBugMetrics(major, minor, "yuchao.li@pingcap.com")
	assert.Nil(t, err)
}

func TestRefreshSprintMetaInfo(t *testing.T) {
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	config := configs.NewConfig(
		"../../"+configs.TestConfig,
		"../../"+configs.TestSecretConfig)
	store.NewStore(config)

	targetSprint := []string{
		"6.5", "6.4", "6.3", "6.2", "6.1", "6.0",
		"5.4",
		"5.3",
	}
	sort.Strings(targetSprint)

	repos, err := store.SelectRepo(&entity.RepoOption{})
	assert.Nil(t, err)

	for _, sprint := range targetSprint {
		major, _ := strconv.Atoi(strings.Split(sprint, ".")[0])
		minor, _ := strconv.Atoi(strings.Split(sprint, ".")[1])
		for _, repo := range *repos {
			err := refreshSprint(major, minor, repo.Owner, repo.Repo)
			if err != nil {
				fmt.Printf("%v \n", err)
			}
		}
	}
}
