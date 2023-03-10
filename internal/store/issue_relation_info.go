package store

import (
	"fmt"

	"tirelease/commons/utils"
	"tirelease/internal/entity"

	"github.com/pkg/errors"
)

// Fill the issues and it's affection infos of IssueRelations
func SelectIssueRelationInfoByJoin(option *entity.IssueRelationInfoOption) (*[]entity.IssueRelationInfoByJoin, error) {
	sql := "select issue.issue_id, group_concat(issue_affect.id) as issue_affect_ids from "
	sql += composeIssueRelationInfoByJoin(option, true)

	// 查询
	var issueRelationInfoJoin []entity.IssueRelationInfoByJoin
	if err := storeGlobalDB.RawWrapper(sql, option).Find(&issueRelationInfoJoin).Error; err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("select issue_relation by raw by join failed, option: %+v", option))
	}

	return &issueRelationInfoJoin, nil
}

func CountIssueRelationInfoByJoin(option *entity.IssueRelationInfoOption) (int64, error) {
	sql := "select count(*) from "
	sql += composeIssueRelationInfoByJoin(option, false)

	// 查询
	var count int64
	if err := storeGlobalDB.RawWrapper(sql, option).Count(&count).Error; err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("count issue_relation by raw by join failed, option: %+v", option))
	}
	return count, nil
}

func composeIssueRelationInfoByJoin(option *entity.IssueRelationInfoOption, isLimit bool) string {
	issueAffectOption := &entity.IssueAffectOption{
		AffectVersion: option.AffectVersion,
		AffectResult:  option.AffectResult,
	}
	isAffectFilter := (option.AffectVersion != "" || option.AffectResult != "")

	sql := ""
	sql += " ( "
	sql += "select * from issue where 1=1 " + IssueWhere(&option.IssueOption)
	// Filter with issue option
	if isAffectFilter {
		sql += " and issue_id in ( "
		sql += " SELECT issue_id "
		sql += " FROM issue_affect "
		sql += " WHERE 1=1 " + IssueAffectWhere(issueAffectOption)
		sql += " )"
	}
	if isLimit {
		sql += IssueOrderBy(&option.IssueOption)
		sql += IssueLimit(&option.IssueOption)
	}
	sql += " ) as issue "
	sql += "left join "
	sql += " ( "
	sql += "select * from issue_affect"
	sql += " ) as issue_affect "
	sql += "on issue.issue_id = issue_affect.issue_id "
	sql += "group by issue.issue_id "
	sql += IssueOrderBy(&option.IssueOption)

	return sql
}

// SelectUntriagedIssueRelationInfo function
// Select issues with their affection info with following limit:
//  1. The issue affects some **active** versions that have not been triaged.
func SelectNeedTriageIssueRelationInfo(option *entity.IssueRelationInfoOption) (*[]entity.IssueRelationInfoByJoin, error) {
	sql := "select issue.issue_id, group_concat(issue_affect.id) as issue_affect_ids from "
	sql += composeNeedTriageIssueRelationInfo(option, true)

	// 查询
	var issueRelationInfoJoin []entity.IssueRelationInfoByJoin
	if err := storeGlobalDB.RawWrapper(sql, option).Find(&issueRelationInfoJoin).Error; err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("select unpicked issue_relation info by raw by join failed, option: %+v", option))
	}

	return &issueRelationInfoJoin, nil
}

func CountNeedTriageIssueRelationInfo(option *entity.IssueRelationInfoOption) (int64, error) {
	sql := "select count(*) from "
	sql += composeNeedTriageIssueRelationInfo(option, false)

	// 查询
	var count int64
	if err := storeGlobalDB.RawWrapper(sql, option).Count(&count).Error; err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("count unpicked issue_relation info  by raw by join failed, option: %+v", option))
	}
	return count, nil
}

func composeNeedTriageIssueRelationInfo(option *entity.IssueRelationInfoOption, isLimit bool) string {
	issueAffectOption := &entity.IssueAffectOption{
		AffectVersion: option.AffectVersion,
		AffectResult:  option.AffectResult,
	}
	isAffectFilter := (option.AffectVersion != "" || option.AffectResult != "")

	sql := ""
	sql += " ( "
	sql += "select * from issue where 1=1 " + IssueWhere(&option.IssueOption)
	// Filter with issue option
	if isAffectFilter {
		sql += " and issue_id in ( "
		sql += " SELECT issue_id "
		sql += " FROM issue_affect "
		sql += " WHERE 1=1 " + IssueAffectWhere(issueAffectOption)
		sql += " )"
	}
	// Filter unpicked issues
	sql += fmt.Sprintf(`
        and issue_id in (
            select affect.issue_id
            from (
                select *
                from issue_affect
                where affect_result = "Yes"
                and affect_version in (
                    select substring_index(name,".",2) as minor_version
                     from release_version
                     WHERE status in (%[1]s)
                )
            ) affect
            left join (
                select *
                    , substring_index(version_name,".",2) as minor_version
                from version_triage
                where triage_result not in ("%[2]s")
            ) triage
            on affect.issue_id = triage.issue_id
            and affect.affect_version = triage.minor_version
            where triage.issue_id is null
        )
    `, utils.Join(entity.ActiveVersionStatus, ",", "\""),
		string(entity.VersionTriageResultUnKnown),
	)
	if isLimit {
		sql += IssueOrderBy(&option.IssueOption)
		sql += IssueLimit(&option.IssueOption)
	}
	sql += " ) as issue "
	sql += "left join "
	sql += " ( "
	sql += "select * from issue_affect"
	sql += " ) as issue_affect "
	sql += "on issue.issue_id = issue_affect.issue_id "
	sql += "group by issue.issue_id "
	sql += IssueOrderBy(&option.IssueOption)

	return sql
}
