package service

import (
	"testing"
	"tirelease/commons/feishu"
	"tirelease/commons/git"
	"tirelease/internal/store"
	"tirelease/utils/configs"

	"github.com/stretchr/testify/assert"
)

func TestNotifySprintReleaseNotesExcel(t *testing.T) {
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	config := configs.NewConfig(
		"../../"+configs.TestConfig,
		"../../"+configs.TestSecretConfig)
	store.NewStore(config)
	// FeishuAPP
	feishu.SetFeishuApp(config.FeiShu.AppId, config.FeiShu.AppSecret)

	major := 6
	minor := 5
	err := NotifySprintReleaseNotesExcel(major, minor, "yuchao.li@pingcap.com")
	assert.Nil(t, err)
}
