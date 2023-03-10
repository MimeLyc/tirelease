package store

import (
	"fmt"
	"time"
	"tirelease/internal/entity"

	"github.com/pkg/errors"
	"gorm.io/gorm/clause"
)

func SelectHotfixReleaseInfos(option *entity.HotfixReleaseInfoOptions) (*[]entity.HotfixReleaseInfo, error) {
	sql := "select * from hotfix_release_info where 1=1" + hotfixReleaseWhere(option) + hotfixReleaseOrderBy(option) + hotfixReleaseLimit(option)
	// 查询
	var release []entity.HotfixReleaseInfo
	if err := storeGlobalDB.RawWrapper(sql, option).Find(&release).Error; err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("find release version: %+v failed", option))
	}

	return &release, nil
}

func CreateOrUpdateHotfixReleaseInfo(release *entity.HotfixReleaseInfo) error {
	// 加工
	if release.CreateTime.IsZero() {
		release.CreateTime = time.Now()
	}
	release.UpdateTime = time.Now()

	// 存储
	if err := storeGlobalDB.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&release).Error; err != nil {
		return errors.Wrap(err, fmt.Sprintf("create hotfix: %+v failed", release))
	}
	return nil
}

func hotfixReleaseWhere(option *entity.HotfixReleaseInfoOptions) string {
	sql := ""

	if option.HotfixName != "" {
		sql += " and hotfix_release_info.hotfix_name = @HotfixName"
	}
	if option.HotfixNames != nil {
		sql += " and hotfix_release_info.hotfix_name in @HotfixNames"
	}

	return sql
}

func hotfixReleaseOrderBy(option *entity.HotfixReleaseInfoOptions) string {
	return option.GetOrderByString()
}

func hotfixReleaseLimit(option *entity.HotfixReleaseInfoOptions) string {
	return option.GetLimitString()
}
