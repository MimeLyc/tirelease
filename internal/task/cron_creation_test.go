package task

import (
	"testing"
	"tirelease/internal/entity"
	"tirelease/internal/store"
	"tirelease/utils/configs"
)

func TestCreateCronTask(t *testing.T) {

	config := configs.NewConfig(
		"../../"+configs.TestConfig,
		"../../"+configs.TestSecretConfig)
	store.NewStore(config)

	task := entity.Task{
		Type: entity.TASK_TYPE_REFRESH_EMPLOYEE,
	}

	// Create the same cron task twice to ensure the cron task is created only once
	CreateCronTask(task)

	CreateCronTask(task)
}
