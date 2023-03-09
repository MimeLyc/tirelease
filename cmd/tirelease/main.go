package main

import (
	"tirelease/api"
	"tirelease/commons/configs"
	"tirelease/commons/database"
	"tirelease/commons/feishu"
	"tirelease/commons/git"
	"tirelease/internal/cron"
	"tirelease/internal/store"
	"tirelease/internal/task"
)

func main() {
	// Load config
	configs.LoadConfig("config.yaml")

	// Connect database
	database.Connect(configs.Config)
	store.InitHrEmployeeDB()

	// Github Client (If Needed: V3 & V4)
	git.Connect(configs.Config.Github.AccessToken)
	git.ConnectV4(configs.Config.Github.AccessToken)

	// FeishuAPP
	feishu.SetFeishuApp(configs.Config.Feishu.AppId, configs.Config.Feishu.AppSecret)

	// Start Cron (If Needed)
	cron.InitCron()

	// Start Task Execution
	go task.StartTaskExecution()

	// Start website && REST-API
	router := api.Routers("website/build/")
	router.Run(":8080")
}
