package main

import (
	"tirelease/api"
	"tirelease/commons/feishu"
	"tirelease/commons/git"
	"tirelease/internal/cron"
	"tirelease/internal/store"
	"tirelease/internal/task"
	"tirelease/utils/configs"
)

var defaultConfigPath = "config.yaml"

func main() {
	// Load config
	config := configs.NewConfig(defaultConfigPath)

	// Connect database
	store.New(config)
	store.InitHrEmployeeDB(config.EmployeeDSN)

	// Github Client (If Needed: V3 & V4)
	git.Connect(config.GitHubAccessToken)
	git.ConnectV4(config.GitHubAccessToken)

	// FeishuAPP
	feishu.SetFeishuApp(config.Feishu.AppId, config.Feishu.AppSecret)

	// Start Cron (If Needed)
	cron.InitCron()

	// Start Task Execution
	go task.StartTaskExecution()

	// Start website && REST-API
	router := api.Routers("website/build/")
	router.Run(":8080")
}
