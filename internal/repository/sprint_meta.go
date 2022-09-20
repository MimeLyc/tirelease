package repository

import (
	"fmt"
	"tirelease/commons/database"
	"tirelease/internal/entity"

	"github.com/pkg/errors"
)

func SelectSprintMetaUnique(option *entity.SprintMetaOption) (*entity.SprintMeta, error) {
	sql := "select * from sprint_meta where 1=1" + SprintMetaWhere(option) + option.GetOrderByString() + option.GetLimitString()
	// 查询
	var sprintMetas []entity.SprintMeta
	if err := database.DBConn.RawWrapper(sql, option).Find(&sprintMetas).Error; err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("find sprint meta: %+v failed", option))
	}

	if len(sprintMetas) == 0 {
		return nil, nil
	}

	return &sprintMetas[0], nil
}

func SprintMetaWhere(option *entity.SprintMetaOption) string {
	sql := ""

	if option.ID != nil {
		sql += " and sprint_meta.id = @ID"
	}
	if option.MinorVersionName != nil {
		sql += " and sprint_meta.minor_version_name = @MinorVersionName"
	}
	if option.Major != nil {
		sql += " and sprint_meta.major = @Major"
	}
	if option.Minor != nil {
		sql += " and sprint_meta.minor = @Minor"
	}
	if option.Repo != nil {
		sql += " and sprint_meta.repo_id = @Repo"
	}
	if option.StartTime != nil {
		sql += " and sprint_meta.start_time = @StartTime"
	}
	if option.CheckoutCommitTime != nil {
		sql += " and sprint_meta.checkout_commit_time = @CheckoutCommitTime"
	}

	return sql
}
