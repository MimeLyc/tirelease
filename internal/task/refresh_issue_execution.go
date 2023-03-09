package task

import (
	"tirelease/internal/entity"
	"tirelease/internal/service"
	"tirelease/internal/store"
)

type RefreshIssueTask struct {
	TaskExecutionBase
}

func (refreshTask RefreshIssueTask) getTaskType() entity.TaskType {
	return entity.TASK_TYPE_REFRESH_ISSUE
}

func (refreshTask RefreshIssueTask) process(task *entity.Task) error {
	repos, err := store.SelectRepo(&entity.RepoOption{})
	if err != nil {
		return err
	}
	releaseVersions, err := store.SelectReleaseVersion(&entity.ReleaseVersionOption{})
	if err != nil {
		return err
	}
	params := &service.RefreshIssueParams{
		Repos:           repos,
		BeforeHours:     -2,
		Batch:           20,
		Total:           500,
		IsHistory:       true,
		ReleaseVersions: releaseVersions,
		Order:           "DESC",
	}

	err = service.CronRefreshIssuesV4(params)

	return err
}

func NewRefreshIssueTask() RefreshIssueTask {
	task := &RefreshIssueTask{}
	task.ITaskExecution = interface{}(task).(ITaskExecution)

	return *task
}
