package entity

import (
	"time"

	"github.com/google/go-github/v41/github"
)

// Struct of Pull Request
type PullRequest struct {
	// DataBase columns
	ID            int64  `json:"id,omitempty"`
	PullRequestID string `json:"pull_request_id,omitempty"`
	Number        int    `json:"number,omitempty"`
	State         string `json:"state,omitempty"`
	Title         string `json:"title,omitempty"`
	Owner         string `json:"owner,omitempty"`
	Repo          string `json:"repo,omitempty"`
	HTMLURL       string `json:"html_url,omitempty"`
	BaseBranch    string `json:"base_branch,omitempty"`

	CreateTime time.Time  `json:"create_time,omitempty"`
	UpdateTime time.Time  `json:"update_time,omitempty"`
	CloseTime  *time.Time `json:"close_time,omitempty"`
	MergeTime  *time.Time `json:"merge_time,omitempty"`

	Merged             bool    `json:"merged,omitempty"`
	MergeableState     *string `json:"mergeable_state,omitempty"`
	CherryPickApproved bool    `json:"cherry_pick_approved,omitempty"`
	AlreadyReviewed    bool    `json:"already_reviewed,omitempty"`

	SourcePullRequestID string `json:"source_pull_request_id,omitempty"`

	LabelsString             string `json:"labels_string,omitempty"`
	AssigneesString          string `json:"assignees_string,omitempty"`
	RequestedReviewersString string `json:"requested_reviewers_string,omitempty"`
	IsReleaseNoteConfirmed   bool   `json:"is_release_note_confirmed,omitempty"`
	ReleaseNote              string `json:"releaseNote,omitempty"`

	// OutPut-Serial
	Labels             *[]github.Label `json:"labels,omitempty" gorm:"-"`
	Assignees          *[]github.User  `json:"assignees,omitempty" gorm:"-"`
	RequestedReviewers *[]github.User  `json:"requested_reviewers,omitempty" gorm:"-"`
	Body               string          `json:"body,omitempty" gorm:"-"`
}

// List Option
type PullRequestOption struct {
	ID                  int64  `json:"id" form:"id"`
	PullRequestID       string `json:"pull_request_id,omitempty" form:"pull_request_id"`
	Number              int    `json:"number,omitempty" form:"number"`
	State               string `json:"state,omitempty" form:"state"`
	Owner               string `json:"owner,omitempty" form:"owner"`
	Repo                string `json:"repo,omitempty" form:"repo"`
	BaseBranch          string `json:"base_branch,omitempty" form:"base_branch"`
	SourcePullRequestID string `json:"source_pull_request_id,omitempty" form:"source_pull_request_id"`
	Merged              *bool  `json:"merged,omitempty"`
	MergeableState      string `json:"mergeable_state,omitempty"`
	CherryPickApproved  *bool  `json:"cherry_pick_approved,omitempty"`
	AlreadyReviewed     *bool  `json:"already_reviewed,omitempty"`

	PullRequestIDs []string `json:"pull_request_ids,omitempty" form:"pull_request_ids"`

	ListOption
}

// DB-Table
func (PullRequest) TableName() string {
	return "pull_request"
}

/**

CREATE TABLE IF NOT EXISTS pull_request (
	id INT(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
	pull_request_id VARCHAR(255) COMMENT 'Pr全局ID',
	number INT(11) NOT NULL COMMENT '当前库ID',
	state VARCHAR(32) NOT NULL COMMENT '状态',
	title VARCHAR(1024) COMMENT '标题',

	owner VARCHAR(255) COMMENT '仓库所有者',
	repo VARCHAR(255) COMMENT '仓库名称',
	html_url VARCHAR(1024) COMMENT '链接',
	base_branch VARCHAR(255) COMMENT '目标分支',

	close_time TIMESTAMP COMMENT '关闭时间',
	create_time TIMESTAMP COMMENT '创建时间',
	update_time TIMESTAMP COMMENT '更新时间',
	merge_time TIMESTAMP COMMENT '合入时间',

	merged BOOLEAN COMMENT '是否已合入',
	mergeable_state VARCHAR(32) COMMENT '可合入状态',
	cherry_pick_approved BOOLEAN COMMENT '是否已进入版本',
	already_reviewed BOOLEAN COMMENT '是否已代码评审',

	source_pull_request_id VARCHAR(255) COMMENT '来源ID',
	labels_string TEXT COMMENT '标签',
	assignees_string TEXT COMMENT '处理人列表',
	requested_reviewers_string TEXT COMMENT '处理人列表',

	PRIMARY KEY (id),
	UNIQUE KEY uk_prid (pull_request_id),
	INDEX idx_state (state),
	INDEX idx_owner_repo (owner, repo),
	INDEX idx_createtime (create_time),
	INDEX idx_sourceprid (source_pull_request_id)
)
ENGINE = INNODB DEFAULT CHARSET = utf8 COMMENT 'pull_request信息表';

**/
