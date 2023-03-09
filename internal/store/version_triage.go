package store

import (
	"fmt"
	"time"

	"tirelease/internal/entity"

	"github.com/pkg/errors"
	"gorm.io/gorm/clause"
)

func CreateVersionTriage(versionTriage *entity.VersionTriage) error {
	if versionTriage.CreateTime.IsZero() {
		versionTriage.CreateTime = time.Now()
	}
	if versionTriage.UpdateTime.IsZero() {
		versionTriage.UpdateTime = time.Now()
	}
	// 存储
	if err := tempDB.DB.Clauses(
		clause.OnConflict{DoNothing: true}).Create(&versionTriage).Error; err != nil {
		return errors.Wrap(err, fmt.Sprintf("create version triage: %+v failed", versionTriage))
	}
	return nil
}

func SelectVersionTriage(option *entity.VersionTriageOption) (*[]entity.VersionTriage, error) {
	sql := "select * from version_triage where 1=1"
	sql += VersionTriageWhere(option) + VersionTriageOrderBy(option) + VersionTriageLimit(option)

	// 查询
	var versionTriages []entity.VersionTriage
	if err := tempDB.RawWrapper(sql, option).Find(&versionTriages).Error; err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("find version triage: %+v failed", option))
	}
	return &versionTriages, nil
}

func SelectVersionTriageUnique(option *entity.VersionTriageOption) (*entity.VersionTriage, error) {
	// 查询
	triages, err := SelectVersionTriage(option)
	if err != nil {
		return nil, err
	}

	// 校验
	if len(*triages) == 0 {
		return nil, nil
	}

	return &((*triages)[0]), nil
}

func CountVersionTriage(option *entity.VersionTriageOption) (int64, error) {
	sql := "select count(*) from version_triage where 1=1"
	sql += VersionTriageWhere(option)

	// 查询
	var count int64
	if err := tempDB.RawWrapper(sql, option).Count(&count).Error; err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("count version triage: %+v failed", option))
	}
	return count, nil
}

func UpdateVersionTriage(versionTriage *entity.VersionTriage) error {
	// 更新
	if versionTriage.UpdateTime.IsZero() {
		versionTriage.UpdateTime = time.Now()
	}
	if err := tempDB.DB.Omit("CreateTime").Save(&versionTriage).Error; err != nil {
		return errors.Wrap(err, fmt.Sprintf("update version triage: %+v failed", versionTriage))
	}
	return nil
}

func CreateOrUpdateVersionTriage(versionTriage *entity.VersionTriage, updatedVars ...entity.VersionTriageUpdatedVar) error {
	// 存储
	if versionTriage.CreateTime.IsZero() {
		versionTriage.CreateTime = time.Now()
	}
	versionTriage.UpdateTime = time.Now()

	if err := tempDB.DB.Clauses(clause.OnConflict{
		DoUpdates: clause.AssignmentColumns(composeUpdatedVars(updatedVars...)),
	}).Create(&versionTriage).Error; err != nil {
		return errors.Wrap(err, fmt.Sprintf("create or update version triage: %+v failed", versionTriage))
	}
	return nil
}

func VersionTriageWhere(option *entity.VersionTriageOption) string {
	sql := ""

	if option.ID != 0 {
		sql += " and version_triage.id = @ID"
	}
	if option.IssueID != "" {
		sql += " and version_triage.issue_id = @IssueID"
	}
	if option.VersionName != "" {
		sql += " and version_triage.version_name = @VersionName"
	}
	if option.TriageResult != "" {
		sql += " and version_triage.triage_result = @TriageResult"
	}
	if option.IssueIDs != nil && len(option.IssueIDs) > 0 {
		sql += " and version_triage.issue_id in @IssueIDs"
	}
	if option.VersionNameList != nil && len(option.VersionNameList) > 0 {
		sql += " and version_triage.version_name in @VersionNameList"
	}

	return sql
}

func VersionTriageOrderBy(option *entity.VersionTriageOption) string {
	sql := ""

	if option.OrderBy != "" {
		sql += " order by " + option.OrderBy
	}
	if option.Order != "" {
		sql += " " + option.Order
	}

	return sql
}

func VersionTriageLimit(option *entity.VersionTriageOption) string {
	sql := ""

	if option.Page != 0 && option.PerPage != 0 {
		option.ListOption.CalcOffset()
		sql += " limit @Offset,@PerPage"
	}

	return sql
}

func composeUpdatedVars(updatedVars ...entity.VersionTriageUpdatedVar) []string {
	if len(updatedVars) == 0 {
		return []string{"update_time", "version_name", "issue_id", "triage_owner", "triage_result", "block_version_release", "due_time", "comment"}
	}

	vars := []string{"update_time"}
	for _, v := range updatedVars {
		if string(v) == "update_time" {
			continue
		}
		vars = append(vars, string(v))
	}
	return vars
}

func DeleteVersionTriagesByIssueId(issueId string) ([]entity.VersionTriage, error) {
	where := fmt.Sprintf("issue_id = '%s'", issueId)
	var triages []entity.VersionTriage

	if err := tempDB.DB.Clauses(clause.Returning{}).Where(where).Delete(&triages).Error; err != nil {
		return nil, err
	}
	return triages, nil
}
