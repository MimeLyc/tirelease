package model

import (
	"fmt"
	"strings"
	"tirelease/internal/entity"
	"tirelease/internal/service/component"
)

func ComponentString(pr PullRequest, issue *entity.Issue) string {
	if issue == nil {
		return pr.Repo
	}

	componentsString := make([]string, 0)
	components := component.GetComponents(issue.Owner, issue.Repo, issue.LabelsString)
	componentsString = append(componentsString, component.ParseToString(components)...)

	if len(componentsString) == 0 {
		return pr.Repo
	}

	return strings.Join(componentsString, ", ")
}

func PullRequestLabelsString(pr PullRequest) string {
	labels := pr.Labels
	labelStrings := make([]string, 0)
	for _, label := range *labels {
		labelStrings = append(labelStrings, *label.Name)
	}

	return strings.Join(labelStrings, ", ")
}

func PullRequestAssigneesString(pr PullRequest) string {
	assignees := pr.Assignees
	assigneeStrings := make([]string, 0)
	for _, assignee := range *assignees {
		assigneeStrings = append(assigneeStrings, *assignee.Login)
	}

	return strings.Join(assigneeStrings, ", ")
}

const GithubUserTemplate = "@[%s](https://github.com/%s)"

func FormatedPullRequestAuthorString(pr PullRequest) string {
	formatedString := fmt.Sprintf(GithubUserTemplate, pr.AuthorGhLogin, pr.AuthorGhLogin)
	return formatedString
}

func FormatedPullRequestAssigneesString(pr PullRequest) string {
	assignees := pr.Assignees
	assigneeStrings := make([]string, 0)
	for _, assignee := range *assignees {
		formatedString := fmt.Sprintf(GithubUserTemplate, *assignee.Login, *assignee.Login)
		assigneeStrings = append(assigneeStrings, formatedString)
	}

	return strings.Join(assigneeStrings, ", ")
}

func IssueAssigneesString(issue entity.Issue) string {
	assignees := issue.Assignees
	assigneeStrings := make([]string, 0)
	for _, assignee := range *assignees {
		assigneeStrings = append(assigneeStrings, *assignee.Login)
	}

	return strings.Join(assigneeStrings, ", ")
}

func FormatedIssueAssigneesString(issue entity.Issue) string {
	assignees := issue.Assignees
	assigneeStrings := make([]string, 0)
	for _, assignee := range *assignees {
		formatedString := fmt.Sprintf(GithubUserTemplate, *assignee.Login, *assignee.Login)
		assigneeStrings = append(assigneeStrings, formatedString)
	}

	return strings.Join(assigneeStrings, ", ")
}

const MdFormatedUrl = "[#%d](%s)"

func FormatedPullrequestUrl(pr PullRequest) string {
	return fmt.Sprintf(MdFormatedUrl, pr.Number, pr.HTMLURL)
}

func FormatedIssueUrl(issue entity.Issue) string {
	return fmt.Sprintf(MdFormatedUrl, issue.Number, issue.HTMLURL)
}

const FormatedReleaseNoteTemplate = "%s %s %s"

func FormatedReleaseNote(pr PullRequest, issue *entity.Issue) string {

	assignees := ""
	if issue != nil {
		assignees = FormatedIssueAssigneesString(*issue)
	}
	if len(assignees) == 0 {
		assignees = FormatedPullRequestAssigneesString(pr)
	}
	if len(assignees) == 0 {
		assignees = FormatedPullRequestAuthorString(pr)
	}

	url := ""
	if issue != nil {
		url = FormatedIssueUrl(*issue)
	}
	if len(url) == 0 {
		url = FormatedPullrequestUrl(pr)
	}
	return fmt.Sprintf(FormatedReleaseNoteTemplate, pr.ReleaseNote, url, assignees)
}

func IssueLabelsString(issue entity.Issue) string {
	labels := issue.Labels
	labelStrings := make([]string, 0)
	for _, label := range *labels {
		labelStrings = append(labelStrings, *label.Name)
	}

	return strings.Join(labelStrings, ", ")
}
