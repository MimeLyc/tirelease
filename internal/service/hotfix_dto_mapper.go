package service

import (
	"fmt"
	"strings"
	"time"
	"tirelease/commons/git"
	"tirelease/internal/dto"
	"tirelease/internal/entity"
	"tirelease/internal/model"
	"tirelease/internal/store"
)

// TODO:
// Refator: dto should depend on service, not the other way around.
// 1. Change all dto model in service package to model packages.
// 2. Change the mapping function to dto package.
func buildHotfixFromRequest(r dto.HotfixSaveRequest) (model.Hotfix, error) {
	hotfix := r.Hotfix
	// fill hotfix name
	if hotfix.Name == "" && hotfix.BaseVersionName != "" && hotfix.OncallID != "" {
		today := time.Now().Format("20060102")
		hotfix.Name = fmt.Sprintf(
			"%s-%s-%s-%s",
			today,
			hotfix.BaseVersionName,
			strings.ToUpper(hotfix.OncallPrefix),
			hotfix.OncallID,
		)
	}

	releaseInfos := make([]model.HotfixReleaseInfo, 0)
	for _, releaseReq := range r.ReleaseInfos {
		releaseInfo, err := buildHotfixReleaseFromRequest(hotfix.Name, releaseReq)
		if err != nil {
			return hotfix, err
		}
		releaseInfos = append(releaseInfos, releaseInfo)
	}
	hotfix.ReleaseInfos = releaseInfos

	return hotfix, nil
}

// buildHotfixReleaseFromRequest build hotfix release info from request
func buildHotfixReleaseFromRequest(hotfixName string, req dto.HotfixReleaseInfoRequest) (model.HotfixReleaseInfo, error) {
	releaseInfo := req.HotfixReleaseInfo
	releaseInfo.HotfixName = hotfixName
	releaseInfo.AssigneeEmail = req.Assignee.Email

	// fill in repo full name
	repos, err := store.SelectRepo(
		&entity.RepoOption{
			Repo: req.Repo,
		},
	)
	if err != nil {
		return releaseInfo, err
	}
	if len(*repos) != 1 {
		return releaseInfo, fmt.Errorf("repo %s not found", req.Repo)
	}
	repo := (*repos)[0]
	releaseInfo.RepoFullName = repo.FullName

	// fill in issues
	issues := make([]model.Issue, 0)
	for _, issue := range req.Issues {
		issueId, err := getIssueId(issue)
		if err != nil {
			return releaseInfo, err
		}
		issues = append(issues, model.Issue{
			Issue: entity.Issue{
				IssueID: issueId,
				Owner:   issue.Owner,
				Repo:    issue.Repo,
				HTMLURL: issue.HTMLURL,
			},
		})
	}
	releaseInfo.Issues = issues

	// fill in master prs
	prs := make([]model.PullRequest, 0)
	for _, pr := range req.MasterPrs {
		prId, err := getPrId(pr)
		if err != nil {
			return releaseInfo, err
		}
		prs = append(prs, model.PullRequest{
			PullRequest: &entity.PullRequest{
				PullRequestID: prId,
				Owner:         pr.Owner,
				Repo:          pr.Repo,
				HTMLURL:       pr.HTMLURL,
			},
		})
	}
	releaseInfo.MasterPrs = prs

	// fill in branch prs
	prs = make([]model.PullRequest, 0)
	for _, pr := range req.BranchPrs {
		prId, err := getPrId(pr)
		if err != nil {
			return releaseInfo, err
		}
		prs = append(prs, model.PullRequest{
			PullRequest: &entity.PullRequest{
				PullRequestID: prId,
				Owner:         pr.Owner,
				Repo:          pr.Repo,
				HTMLURL:       pr.HTMLURL,
			},
		})
	}
	releaseInfo.BranchPrs = prs

	return releaseInfo, nil
}

func getIssueId(issue dto.HotfixIssue) (string, error) {
	entity, err := store.SelectIssueUnique(
		&entity.IssueOption{
			Owner:  issue.Owner,
			Repo:   issue.Repo,
			Number: issue.Number,
		},
	)
	if _, ok := err.(store.DataNotFoundError); ok {
		issueId, err := RefreshIssueV4ByNumber(
			issue.Owner,
			issue.Repo,
			issue.Number,
		)
		if err != nil {
			return "", err
		}
		return issueId, nil
	} else if err != nil {
		return "", err
	} else {
		return entity.IssueID, nil
	}

}

func getPrId(req dto.HotfixPr) (string, error) {
	entity, err := store.SelectPullRequestUnique(
		&entity.PullRequestOption{
			Owner:  req.Owner,
			Repo:   req.Repo,
			Number: req.Number,
		},
	)
	if _, ok := err.(store.DataNotFoundError); ok {
		prId, err := refreshAndGetPrId(req)
		if err != nil {
			return "", err
		}
		return prId, nil
	} else if err != nil {
		return "", err
	} else {
		return entity.PullRequestID, nil
	}

}

func refreshAndGetPrId(req dto.HotfixPr) (string, error) {
	pr, _, err := git.Client.GetPullRequestByNumber(req.Owner, req.Repo, req.Number)
	if err != nil {
		return "", err
	}

	err = WebhookRefreshPullRequestV3(pr)
	if err != nil {
		return "", err
	}
	return *pr.NodeID, nil
}
