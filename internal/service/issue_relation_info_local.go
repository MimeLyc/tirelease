package service

import (
	"fmt"
	"strconv"
	"strings"

	"tirelease/commons/git"
	"tirelease/commons/utils"
	"tirelease/internal/dto"
	"tirelease/internal/entity"
	"tirelease/internal/model"
	"tirelease/internal/repository"
	"tirelease/internal/service/component"
)

// ============================================================================
// ============================================================================ CURD Of IssueRelationInfo
// Get relation infomations of target issue
// relation infomations include:
//  1. Issue : Issue basic info
//  2. IssueAffects : The minor versions affected by the issue
//  3. IssuePrRelations : The pull requests related to the issue **regardless** of the version**
//  4. PullRequests	: The pull requests related to the issue **in the version**
//  5. VersionTriages : The version triage history of the issue
//
// ============================================================================
// TODO: Decouple the infos of current version from the infos of all versions
//
//	    meta: Issue
//	    current version infos: PullRequests
//		   all issue info：IssueAffects, IssuePrRelations, VersionTriages
func FindIssueRelationInfo(query *dto.IssueRelationInfoQuery) (*[]dto.IssueRelationInfo, *entity.ListResponse, error) {
	option := query.Map2EntityOption()

	// select issues and affectioninfos
	joins, count, err := FindIssueRelationEntitys(query)
	if err != nil {
		return nil, nil, err
	}

	response := &entity.ListResponse{
		TotalCount: count,
		Page:       option.IssueOption.Page,
		PerPage:    option.IssueOption.PerPage,
	}
	response.CalcTotalPage()

	// Get all issue ids for further batch select of other entities
	issueIDs := getIssueIDs(*joins)

	// Get all affected minor versions of the issue
	issueAffectAll, err := getIssueAffectVersions(*joins)
	if err != nil {
		return nil, nil, err
	}

	issueAll, err := getIssues(issueIDs)
	if err != nil {
		return nil, nil, err
	}
	issueAll = filterIssuesByComponent(issueAll, option.Component)

	// The pull requests related to the issue **regardless** of the version**
	issuePrRelationAll, err := model.SelectIssuePrRelationByIds(issueIDs)
	if err != nil {
		return nil, nil, err
	}

	// Get pullrequests whose base branch **regardless** of the version**
	// option.baseBranch
	pullRequestAll, err := getRelatedPullRequests(issuePrRelationAll)
	if err != nil {
		return nil, nil, err
	}

	// Get pullrequests whose base branch **in the version**
	versionPRs := getSameVersionPullRequests(pullRequestAll, option.BaseBranch)

	versionTriageAll, err := getVersionTriages(issueIDs, option.VersionStatus)
	if err != nil {
		return nil, nil, err
	}
	// Get all version-triage-merge-status histories of the issue
	versionTriageAll = fillVersionTriageMergeStatus(versionTriageAll, pullRequestAll, issuePrRelationAll)

	// compose
	issueRelationInfos := composeIssueRelationInfos(issueAll, issueAffectAll, issuePrRelationAll, versionPRs, versionTriageAll)

	return &issueRelationInfos, response, nil
}

func SaveIssueRelationInfo(issueRelationInfo *dto.IssueRelationInfo) error {

	if issueRelationInfo == nil {
		return nil
	}

	// Save Issue
	if issueRelationInfo.Issue != nil {
		if err := repository.CreateOrUpdateIssue(issueRelationInfo.Issue); nil != err {
			return err
		}
	}

	// Save IssueAffects
	if issueRelationInfo.IssueAffects != nil {
		for _, issueAffect := range *issueRelationInfo.IssueAffects {
			if err := repository.CreateOrUpdateIssueAffect(&issueAffect); nil != err {
				return err
			}
		}
	}

	// Save IssuePrRelations
	if issueRelationInfo.IssuePrRelations != nil {
		for _, issuePrRelation := range *issueRelationInfo.IssuePrRelations {
			if err := repository.CreateIssuePrRelation(&issuePrRelation); nil != err {
				return err
			}
		}
	}

	// Save PullRequests
	if issueRelationInfo.PullRequests != nil {
		for _, pullRequest := range *issueRelationInfo.PullRequests {
			if err := repository.CreateOrUpdatePullRequest(&pullRequest); nil != err {
				return err
			}
		}
	}

	return nil
}

func composeIssueRelationInfos(issueAll []entity.Issue, issueAffectAll []entity.IssueAffect,
	issuePrRelationAll []entity.IssuePrRelation, pullRequestAll []entity.PullRequest,
	versionTriageAll []entity.VersionTriage) []dto.IssueRelationInfo {

	// compose
	issueRelationInfos := make([]dto.IssueRelationInfo, 0)
	for index := range issueAll {
		issue := issueAll[index]

		issueRelationInfo := &dto.IssueRelationInfo{}
		issueRelationInfo.Issue = &issue

		issueAffects := make([]entity.IssueAffect, 0)
		if len(issueAffectAll) > 0 {
			for i := range issueAffectAll {
				issueAffect := issueAffectAll[i]
				if issueAffect.IssueID == issue.IssueID {
					issueAffects = append(issueAffects, issueAffect)
				}
			}
		}
		issueRelationInfo.IssueAffects = &issueAffects

		issuePrRelations := make([]entity.IssuePrRelation, 0)
		pullRequests := make([]entity.PullRequest, 0)
		if len(issuePrRelationAll) > 0 {
			for i := range issuePrRelationAll {
				issuePrRelation := issuePrRelationAll[i]
				if issuePrRelation.IssueID != issue.IssueID {
					continue
				}

				issuePrRelations = append(issuePrRelations, issuePrRelation)
				if len(pullRequestAll) > 0 {
					for j := range pullRequestAll {
						pullRequest := pullRequestAll[j]
						if pullRequest.PullRequestID == issuePrRelation.PullRequestID {
							pullRequests = append(pullRequests, pullRequest)
						}
					}
				}
			}
		}
		issueRelationInfo.IssuePrRelations = &issuePrRelations
		issueRelationInfo.PullRequests = &pullRequests

		versionTriages := make([]entity.VersionTriage, 0)
		if len(versionTriageAll) > 0 {
			for i := range versionTriageAll {
				versionTriage := versionTriageAll[i]
				if versionTriage.IssueID == issue.IssueID {
					versionTriages = append(versionTriages, versionTriage)
				}
			}
		}
		issueRelationInfo.VersionTriages = &versionTriages

		issueRelationInfos = append(issueRelationInfos, *issueRelationInfo)
	}

	return issueRelationInfos
}

func getIssueIDs(joins []entity.IssueRelationInfoByJoin) []string {
	issueIDs := make([]string, 0)
	for i := range joins {
		join := joins[i]
		issueIDs = append(issueIDs, join.IssueID)
	}

	return issueIDs
}

func getIssueAffectVersions(joins []entity.IssueRelationInfoByJoin) ([]entity.IssueAffect, error) {
	issueAffectIDs := make([]int64, 0)
	for i := range joins {
		join := (joins)[i]
		ids := strings.Split(join.IssueAffectIDs, ",")
		for _, id := range ids {
			idint, _ := strconv.Atoi(id)
			issueAffectIDs = append(issueAffectIDs, int64(idint))
		}
	}

	issueAffectAll := make([]entity.IssueAffect, 0)

	if len(issueAffectIDs) > 0 {
		issueAffectOption := &entity.IssueAffectOption{
			IDs: issueAffectIDs,
		}
		issueAffectAlls, err := repository.SelectIssueAffect(issueAffectOption)
		if nil != err {
			return nil, err
		}
		issueAffectAll = append(issueAffectAll, (*issueAffectAlls)...)
	}

	return issueAffectAll, nil
}

func getIssues(issueIDs []string) ([]entity.Issue, error) {
	issueAll := make([]entity.Issue, 0)
	if len(issueIDs) > 0 {
		issueOption := &entity.IssueOption{
			IssueIDs: issueIDs,
		}
		issueAlls, err := model.SelectIssues(issueOption)
		if nil != err {
			return nil, err
		}
		issueAll = append(issueAll, (*issueAlls)...)
	}

	return issueAll, nil
}

func getRelatedPullRequests(relatedPrs []entity.IssuePrRelation) ([]entity.PullRequest, error) {
	pullRequestIDs := make([]string, 0)
	pullRequestAll := make([]entity.PullRequest, 0)

	if len(relatedPrs) > 0 {
		for i := range relatedPrs {
			issuePrRelation := relatedPrs[i]
			pullRequestIDs = append(pullRequestIDs, issuePrRelation.PullRequestID)
		}
		pullRequestOption := &entity.PullRequestOption{
			PullRequestIDs: pullRequestIDs,
		}
		pullRequestAlls, err := repository.SelectPullRequest(pullRequestOption)
		if nil != err {
			return nil, err
		}
		pullRequestAll = append(pullRequestAll, (*pullRequestAlls)...)
	}

	return pullRequestAll, nil
}

func getSameVersionPullRequests(pullRequestAll []entity.PullRequest, baseBranch string) []entity.PullRequest {
	if baseBranch == "" {
		return pullRequestAll
	}

	pullRequests := make([]entity.PullRequest, 0)

	if len(pullRequestAll) == 0 {
		return pullRequests
	}

	for i := range pullRequestAll {
		pullRequest := pullRequestAll[i]
		if pullRequest.BaseBranch == baseBranch {
			pullRequests = append(pullRequests, pullRequest)
		}
	}

	return pullRequests
}

func getVersionTriages(issueIDs []string, versionStatus entity.ReleaseVersionStatus) ([]entity.VersionTriage, error) {
	versionTriageAll := make([]entity.VersionTriage, 0)
	if len(issueIDs) > 0 {
		versionTriageOption := &entity.VersionTriageOption{
			IssueIDs: issueIDs,
		}
		versionTriageAlls, err := repository.SelectVersionTriage(versionTriageOption)
		if nil != err {
			return nil, err
		}

		versionTriageAll = append(versionTriageAll, (*versionTriageAlls)...)
	}

	return versionTriageAll, nil
}

func fillVersionTriageMergeStatus(versionTriages []entity.VersionTriage, pullrequestAll []entity.PullRequest, issuePrRelations []entity.IssuePrRelation) []entity.VersionTriage {
	triages := make([]entity.VersionTriage, 0)
	if len(versionTriages) == 0 {
		return triages
	}

	for i := range versionTriages {
		versionTriage := versionTriages[i]
		prIDs := make([]string, 0)

		for _, relation := range issuePrRelations {
			if relation.IssueID == versionTriage.IssueID {
				prIDs = append(prIDs, relation.PullRequestID)
			}
		}

		major, minor, _, _ := ComposeVersionAtom(versionTriage.VersionName)
		minorVersion := fmt.Sprintf("%d.%d", major, minor)
		baseBranch := git.ReleaseBranchPrefix + minorVersion

		prs := make([]entity.PullRequest, 0)
		for _, pr := range pullrequestAll {
			if pr.BaseBranch != baseBranch {
				continue
			}

			for _, prID := range prIDs {
				if prID == pr.PullRequestID {
					prs = append(prs, pr)
				}
			}
		}

		versionTriage.MergeStatus = ComposeVersionTriageMergeStatus(prs)
		triages = append(triages, versionTriage)
	}

	return triages
}

func filterIssuesByComponent(issues []entity.Issue, filterComponent component.Component) []entity.Issue {
	if filterComponent == component.Component("") {
		return issues
	}

	result := make([]entity.Issue, 0)
	for _, issue := range issues {
		issue := issue
		issueComponents := component.GetComponents(issue.Owner, issue.Repo, issue.LabelsString)
		if utils.Contains(issueComponents, filterComponent) {
			result = append(result, issue)
		}
	}

	return result
}

func MapToVersionTriageInfo(versionTriage model.IssueVersionTriage) dto.VersionTriageInfo {
	triageEntity := versionTriage.MapToEntity()
	return dto.VersionTriageInfo{
		ReleaseVersion: versionTriage.Version.ReleaseVersion,
		IsFrozen:       versionTriage.Version.IsFrozen(),
		IsAccept:       versionTriage.PickTriage.IsAccept(),

		VersionTriage:            &triageEntity,
		VersionTriageMergeStatus: versionTriage.GetMergeStatus(),
		// deprecated: IssueRelationInfo in the related API is not used.
		IssueRelationInfo: &dto.IssueRelationInfo{
			Issue: versionTriage.Issue,
			IssueAffects: &[]entity.IssueAffect{
				{
					IssueID:       versionTriage.Issue.IssueID,
					AffectVersion: versionTriage.Version.ComposeVersionMinorName(),
					AffectResult:  versionTriage.Affect,
				},
			},
			IssuePrRelations: nil,
			PullRequests:     &versionTriage.RelatedPrs,
			VersionTriages:   versionTriage.HistoricalTriages,
		},
	}
}

func FindIssueRelationEntitys(query *dto.IssueRelationInfoQuery) (*[]entity.IssueRelationInfoByJoin, int64, error) {
	option := query.Map2EntityOption()

	if query.IsNeedTriage {
		joins, err := repository.SelectNeedTriageIssueRelationInfo(option)
		if nil != err {
			return nil, 0, err
		}

		count, err := repository.CountNeedTriageIssueRelationInfo(option)
		if nil != err {
			return nil, 0, err
		}

		return joins, count, nil
	} else {
		joins, err := repository.SelectIssueRelationInfoByJoin(option)
		if nil != err {
			return nil, 0, err
		}

		count, err := repository.CountIssueRelationInfoByJoin(option)
		if nil != err {
			return nil, 0, err
		}

		return joins, count, nil
	}
}
