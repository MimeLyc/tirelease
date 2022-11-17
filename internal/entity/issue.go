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
	IssueID       string              `json:"issue_id,omitempty" form:"issue_id"`
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

/**

CREATE TABLE IF NOT EXISTS issue (
	id INT(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
	issue_id VARCHAR(255) NOT NULL COMMENT 'Issue全局ID',
	number INT(11) NOT NULL COMMENT '当前库ID',
	state VARCHAR(32) NOT NULL COMMENT '状态',
	title VARCHAR(1024) COMMENT '标题',
	owner VARCHAR(255) COMMENT '仓库所有者',
	repo VARCHAR(255) COMMENT '仓库名称',
	html_url VARCHAR(1024) COMMENT '链接',

	close_time TIMESTAMP COMMENT '关闭时间',
	create_time TIMESTAMP COMMENT '创建时间',
	update_time TIMESTAMP COMMENT '更新时间',

	labels_string TEXT COMMENT '标签',
	assignees_string TEXT COMMENT '处理人列表',

	closed_by_pull_request_id VARCHAR(255) COMMENT '处理的PR',
	severity_label VARCHAR(255) COMMENT '严重等级',
	type_label VARCHAR(255) COMMENT '类型',

	PRIMARY KEY (id),
	UNIQUE KEY uk_issueid (issue_id),
	INDEX idx_state (state),
	INDEX idx_owner_repo (owner, repo),
	INDEX idx_createtime (create_time),
	INDEX idx_updatetime (update_time),
	INDEX idx_closetime (close_time),
	INDEX idx_severitylabel (severity_label),
	INDEX idx_typelabel (type_label)
)
ENGINE = INNODB DEFAULT CHARSET = utf8 COMMENT 'issue信息表';

**/
