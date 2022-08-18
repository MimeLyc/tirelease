package feishu

import (
	"fmt"
	"testing"
	"tirelease/commons/configs"
)

func TestGetTokenFromApp(t *testing.T) {
	configs.LoadConfig("../../config.yaml")
	config := configs.Config
	fmt.Print(GetAccessTokenFromApp(config.Feishu.AppId, config.Feishu.AppSecret))
}
