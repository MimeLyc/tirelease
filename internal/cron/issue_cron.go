package cron

import (
	"tirelease/commons/cron"
	"tirelease/internal/entity"
	"tirelease/internal/service"
	"tirelease/internal/store"
)

func IssueCron() {
	// Cron 表达式及功能方法
	repos, err := store.SelectRepo(&entity.RepoOption{})
	if err != nil {
		return
	}
	releaseVersions, err := store.SelectReleaseVersion(&entity.ReleaseVersionOption{})
	if err != nil {
		return
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
	cron.Create("0 0 */1 * * ?", func() { service.CronRefreshIssuesV4(params) })
}
