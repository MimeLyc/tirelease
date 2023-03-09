package cron

import (
	"tirelease/commons/cron"
	"tirelease/internal/entity"
	"tirelease/internal/service"
	"tirelease/internal/store"
)

func PullRequestCron() {
	// Cron 表达式及功能方法
	repoOption := &entity.RepoOption{}
	repos, err := store.SelectRepo(repoOption)
	if err != nil {
		return
	}
	params := &service.RefreshPullRequestParams{
		Repos:       repos,
		BeforeHours: -2,
		Batch:       20,
		Total:       500,
	}
	cron.Create("0 0 */1 * * ?", func() { service.CronRefreshPullRequestV4(params) })

	// cron.Create("0 0 */2 * * ?", func() { service.CronMergeRetryPullRequestV3() })
}
