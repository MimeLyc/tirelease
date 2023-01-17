package entity

import (
	"strings"
	"time"

	"tirelease/commons/git"
	"tirelease/internal/service/component"

	"github.com/google/go-github/v41/github"
	"github.com/shurcooL/githubv4"
)

// Struct of Issue
type Issue struct {
	// DataBase Column
	ID      int64  `json:"id,omitempty"`
	IssueID string `json:"issue_id,omitempty" excel:"issue_id"`
	Number  int    `json:"number,omitempty" excel:"issue_number"`
	State   string `json:"state,omitempty" excel:"issue_state"`
	Title   string `json:"title,omitempty" excel:"issue_title"`
	Owner   string `json:"owner,omitempty" excel:"repo_org"`
	Repo    string `json:"repo,omitempty" excel:"repo_name"`
	HTMLURL string `json:"html_url,omitempty" excel:"issue_url"`

	CreateTime time.Time  `json:"create_time,omitempty" excel:"issue_create_time"`
	UpdateTime time.Time  `json:"update_time,omitempty"`
	CloseTime  *time.Time `json:"close_time,omitempty" excel:"issue_close_time"`

	LabelsString    string `json:"labels_string,omitempty" excel:"issue_labels"`
	AssigneesString string `json:"assignees_string,omitempty" excel:"assignees_string"`

	ClosedByPullRequestID string `json:"closed_by_pull_request_id,omitempty"`
	SeverityLabel         string `json:"severity_label,omitempty" excel:"issue_severity"`
	TypeLabel             string `json:"type_label,omitempty" excel:"issue_type"`
	AuthorGHLogin         string `json:"author_gh_login,omitempty" excel:"author_gh_login"`

	// OutPut-Serial
	Labels *[]github.Label `json:"labels,omitempty" gorm:"-"`
	// **deprecated**: refactor after the service level Issue is implemented
	Assignees *[]github.User `json:"assignees,omitempty" gorm:"-"`
	// **deprecated**: refactor after the service level Issue is implemented
	AssignedEmployees *[]Employee           `json:"assigned_employees,omitempty" gorm:"-"`
	Components        []component.Component `json:"components,omitempty" gorm:"-"`
	Author            github.User           `json:"author,omitempty" gorm:"-"`

	// TODO create an version+issue model to maintain the fixed status
	IsFixed bool `excel:"is_fixed"  gorm:"-"`
}

// List Option
type IssueOption struct {
	ID            int64               `json:"id" form:"id"`
	IssueID       string              `json:"issue_id,omitempty" form:"issue_id" uri:"issue_id"`
	Number        int                 `json:"number,omitempty" form:"number"`
	State         string              `json:"state,omitempty" form:"state"`
	Owner         string              `json:"owner,omitempty" form:"owner"`
	Repo          string              `json:"repo,omitempty" form:"repo"`
	Component     component.Component `json:"components,omitempty" form:"components"`
	SeverityLabel string              `json:"severity_label,omitempty" form:"severity_label"`
	TypeLabel     string              `json:"type_label,omitempty" form:"type_label"`

	CreateTime        time.Time `json:"create_time,omitempty"`
	UpdateTime        time.Time `json:"update_time,omitempty"`
	CloseTime         time.Time `json:"close_time,omitempty"`
	IssueIDs          []string  `json:"issue_ids,omitempty" form:"issue_ids"`
	SeverityLabels    []string  `json:"severity_labels,omitempty" form:"severity_labels"`
	NotSeverityLabels []string  `json:"not_severity_labels,omitempty" form:"not_severity_labels"`
	CreateTimeEnd     time.Time `json:"create_time_end,omitempty"`
	CloseTimeEnd      time.Time `json:"close_time_end,omitempty"`

	ListOption
}

// DB-Table
func (Issue) TableName() string {
	return "issue"
}

// ComposeIssueFromV3
// TODO: refactor refer to pull_request_composer.go
func ComposeIssueFromV3(issue *github.Issue) *Issue {
	severityLabel := ""
	typeLabel := ""
	labels := &[]github.Label{}
	for i := range issue.Labels {
		node := issue.Labels[i]
		label := &github.Label{
			Name:  node.Name,
			Color: node.Color,
		}
		*labels = append(*labels, *label)

		if strings.HasPrefix(*label.Name, git.SeverityLabel) {
			severityLabel = *label.Name
		}
		if strings.HasPrefix(*label.Name, git.TypeLabel) {
			typeLabel = *label.Name
		}
	}
	assignees := &[]github.User{}
	for i := range issue.Assignees {
		node := issue.Assignees[i]
		user := &github.User{
			Login: node.Login,
		}
		*assignees = append(*assignees, *user)
	}
	url := strings.Split(*issue.RepositoryURL, "/")
	owner := url[len(url)-2]
	repo := url[len(url)-1]
	author := issue.User

	return &Issue{
		IssueID: *issue.NodeID,
		Number:  *issue.Number,
		State:   strings.ToLower(*issue.State),
		Title:   *issue.Title,
		Owner:   owner,
		Repo:    repo,
		HTMLURL: *issue.HTMLURL,

		CreateTime: *issue.CreatedAt,
		UpdateTime: *issue.UpdatedAt,
		CloseTime:  issue.ClosedAt,

		Labels:    labels,
		Assignees: assignees,

		SeverityLabel: severityLabel,
		TypeLabel:     typeLabel,
		Author:        *author,
		AuthorGHLogin: *author.Login,
	}
}

// ComposeIssueFromV4
// TODO: v4 implement by tony at 2022/02/14
func ComposeIssueFromV4(issueFiled *git.IssueField) *Issue {
	severityLabel := ""
	typeLabel := ""
	labels := &[]github.Label{}
	for i := range issueFiled.Labels.Nodes {
		node := issueFiled.Labels.Nodes[i]
		label := github.Label{
			Name: github.String(string(node.Name)),
		}
		if node.Color != "" {
			label.Color = github.String(string(node.Color))
		}
		*labels = append(*labels, label)

		if strings.HasPrefix(*label.Name, git.SeverityLabel) {
			severityLabel = *label.Name
		}
		if strings.HasPrefix(*label.Name, git.TypeLabel) {
			typeLabel = *label.Name
		}
	}
	assignees := &[]github.User{}
	for i := range issueFiled.Assignees.Nodes {
		node := issueFiled.Assignees.Nodes[i]
		user := github.User{
			Login: (*string)(&node.Login),
		}
		*assignees = append(*assignees, user)
	}
	closedByPrID := ""
	if issueFiled.State == githubv4.IssueStateClosed {
		for _, edge := range issueFiled.TimelineItems.Edges {
			closer := edge.Node.ClosedEvent.Closer.PullRequest
			if closer.Number != 0 {
				closedByPrID = closer.ID.(string)
			}
		}
	}
	author := github.User{
		Login: (*string)(&issueFiled.Author.Login),
	}

	issue := &Issue{
		IssueID: issueFiled.ID.(string),
		Number:  int(issueFiled.Number),
		State:   strings.ToLower(string(issueFiled.State)),
		Title:   string(issueFiled.Title),
		Owner:   string(issueFiled.Repository.Owner.Login),
		Repo:    string(issueFiled.Repository.Name),
		HTMLURL: string(issueFiled.Url),

		CreateTime: issueFiled.CreatedAt.Time,
		UpdateTime: issueFiled.UpdatedAt.Time,

		Labels:    labels,
		Assignees: assignees,

		ClosedByPullRequestID: closedByPrID,
		SeverityLabel:         severityLabel,
		TypeLabel:             typeLabel,
		Author:                author,
		AuthorGHLogin:         *author.Login,
	}
	if issueFiled.ClosedAt != nil {
		issue.CloseTime = &issueFiled.ClosedAt.Time
	}

	return issue
}

// IssueRelationInfo
type IssueRelationInfoOption struct {
	// Issue
	IssueOption

	CreatedAtStamp    int64 `json:"created_at_stamp" form:"created_at_stamp"`
	CreatedAtEndStamp int64 `json:"created_at_stamp_end" form:"created_at_stamp_end"`
	UpdatedAtStamp    int64 `json:"updated_at_stamp" form:"updated_at_stamp"`
	ClosedAtStamp     int64 `json:"closed_at_stamp" form:"closed_at_stamp"`
	ClosedAtEndStamp  int64 `json:"closed_at_stamp_end" form:"closed_at_stamp_end"`

	// Filter Option
	AffectVersion string               `json:"affect_version,omitempty" form:"affect_version" uri:"affect_version"`
	AffectResult  AffectResultResult   `json:"affect_result,omitempty" form:"affect_result" uri:"affect_result"`
	BaseBranch    string               `json:"base_branch,omitempty" form:"base_branch" uri:"base_branch"`
	VersionStatus ReleaseVersionStatus `json:"version_status,omitempty" form:"version_status" uri:"version_status"`
}

// Join IssueRelationInfo
type IssueRelationInfoByJoin struct {
	// issue
	IssueID string `json:"issue_id,omitempty"`

	// issue_affect
	IssueAffectIDs string `json:"issue_affect_ids,omitempty"`
}
