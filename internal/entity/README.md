-- Issue Table Creation DML

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
