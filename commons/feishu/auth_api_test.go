package feishu

import (
	"fmt"
	"testing"
	"tirelease/utils/configs"
)

func TestGetTokenFromApp(t *testing.T) {
	config := configs.NewConfig(
		"../../"+configs.TestConfig,
		"../../"+configs.TestSecretConfig)
	fmt.Print(GetAccessTokenFromApp(config.FeiShu.AppId, config.FeiShu.AppSecret))
}
