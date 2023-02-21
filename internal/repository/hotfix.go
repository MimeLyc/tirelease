package repository

import (
	"fmt"
	"time"
	"tirelease/commons/database"
	"tirelease/internal/entity"

	"github.com/pkg/errors"
	"gorm.io/gorm/clause"
)

func SelectHotfixes(option *entity.HotfixOptions) (*[]entity.Hotfix, error) {
	sql := "select * from hotfix where 1=1" + hotfixWhere(option) + hotfixOrderBy(option) + hotfixLimit(option)
	// 查询
	var hotfix []entity.Hotfix
	if err := database.DBConn.RawWrapper(sql, option).Find(&hotfix).Error; err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("find release version: %+v failed", option))
	}

	return &hotfix, nil
}

func SelectFirstHotfixes(option *entity.HotfixOptions) (*entity.Hotfix, error) {
	sql := "select * from hotfix where 1=1" + hotfixWhere(option) + hotfixOrderBy(option) + hotfixLimit(option)
	// 查询
	var hotfix []entity.Hotfix
	if err := database.DBConn.RawWrapper(sql, option).Find(&hotfix).Error; err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("find release version: %+v failed", option))
	}

	if len(hotfix) == 0 {
		return nil, DataNotFoundError{}
	}

	return &hotfix[0], nil
}

func CreateOrUpdateHotfix(hotfix *entity.Hotfix) error {
	// 加工
	if hotfix.CreateTime.IsZero() {
		hotfix.CreateTime = time.Now()
	}
	hotfix.UpdateTime = time.Now()

	// 存储
	if err := database.DBConn.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&hotfix).Error; err != nil {
		return errors.Wrap(err, fmt.Sprintf("create hotfix: %+v failed", hotfix))
	}
	return nil
}

func hotfixWhere(option *entity.HotfixOptions) string {
	sql := ""

	if option.ID != 0 {
		sql += " and hotfix.id = @ID"
	}
	if option.Name != "" {
		sql += " and hotfix.name = @Name"
	}
	if option.BaseVersionName != "" {
		sql += " and hotfix.base_version = @BaseVersionName"
	}
	if option.CreatorHrId != "" {
		sql += " and hotfix.creator_hr_id = @CreatorHrId"
	}
	if option.Status != "" {
		sql += " and hotfix.status = @Status"
	}
	if !option.IsDeleted {
		sql += " and hotfix.is_deleted = false"
	}

	return sql
}

func hotfixOrderBy(option *entity.HotfixOptions) string {
	return option.GetOrderByString()
}

func hotfixLimit(option *entity.HotfixOptions) string {
	return option.GetLimitString()
}
