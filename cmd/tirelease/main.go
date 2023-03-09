package main

import (
	"flag"
	"tirelease/api"
	"tirelease/commons/feishu"
	"tirelease/commons/git"
	"tirelease/internal/cron"
	"tirelease/internal/task"
	"tirelease/utils/configs"
)

var (
	version   = "latest"
	buildTime = "none"

	configPath = flag.String("config", configs.TestConfig, "specify the config path")
	secretDir  = flag.String("secret", configs.TestSecretConfig, "specify the secret config directory")
)

func main() {
	flag.Parse()
	config := configs.NewConfig(*configPath, *secretDir)

	// Connect database
	store.NewStore(config)

	// Github Client (If Needed: V3 & V4)
	git.Connect(config.GitHubAccessToken)
	git.ConnectV4(config.GitHubAccessToken)

	// FeishuAPP
	feishu.SetFeishuApp(config.FeiShu.AppId, config.FeiShu.AppSecret)

	// Start Cron (If Needed)
	cron.InitCron()

	// Start Task Execution
	go task.StartTaskExecution()

	// Start website && REST-API
	router := api.Routers("website/build/")
	router.Run(":8080")
}
