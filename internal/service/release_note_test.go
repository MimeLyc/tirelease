package service

import (
	"testing"
	"tirelease/commons/database"
	"tirelease/commons/feishu"
	"tirelease/commons/git"
	"tirelease/utils/configs"

	"github.com/stretchr/testify/assert"
)

func TestNotifySprintReleaseNotesExcel(t *testing.T) {
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	database.Connect(generateConfig())
	configs.LoadConfig("../../config.yaml")
	// FeishuAPP
	feishu.SetFeishuApp(configs.Config.Feishu.AppId, configs.Config.Feishu.AppSecret)

	major := 6
	minor := 5
	err := NotifySprintReleaseNotesExcel(major, minor, "yuchao.li@pingcap.com")
	assert.Nil(t, err)
}
