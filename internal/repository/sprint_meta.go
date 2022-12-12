package repository

import (
	"fmt"
	"tirelease/commons/database"
	"tirelease/internal/entity"

	"github.com/pkg/errors"
	"gorm.io/gorm/clause"
)

func CreateOrUpdateSprint(sprint *entity.SprintMeta) error {
	// 存储
	if err := database.DBConn.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&sprint).Error; err != nil {
		return errors.Wrap(err, fmt.Sprintf("create or update sprint: %+v failed", sprint))
	}

	return nil
}

func SelectSprintMetaUnique(option *entity.SprintMetaOption) (*entity.SprintMeta, error) {
	sql := "select * from sprint_info where 1=1" + SprintMetaWhere(option) + option.GetOrderByString() + option.GetLimitString()
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
		sql += " and sprint_info.id = @ID"
	}
	if option.MinorVersionName != nil {
		sql += " and sprint_info.minor_version_name = @MinorVersionName"
	}
	if option.Major != nil {
		sql += " and sprint_info.major = @Major"
	}
	if option.Minor != nil {
		sql += " and sprint_info.minor = @Minor"
	}
	if option.Repo != nil {
		repoFK := option.Repo.FullName
		sql += fmt.Sprintf(" and sprint_info.repo_full_name = \"%s\"", repoFK)
	}
	if option.StartTime != nil {
		sql += " and sprint_info.start_time = @StartTime"
	}
	if option.CheckoutCommitTime != nil {
		sql += " and sprint_info.checkout_commit_time = @CheckoutCommitTime"
	}

	return sql
}
