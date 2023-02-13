package feishu_handler

import (
	"errors"
	"strconv"
	"tirelease/commons/git"
	"tirelease/internal/entity"
	"tirelease/internal/service"
)

func deliverIssueCmd(receive MsgReceiveV1, content content) error {
	switch content.cmd {
	case "approve":
		return approveIssueCmd(content)
	}
	return nil
}

func approveIssueCmd(content content) error {
	owner := content.extractByKey("owner")
	repo := content.extractByKey("repo")
	issueString := content.extractByKey("issue")
	issue, err := strconv.Atoi(issueString)

	if len(owner) == 0 || len(repo) == 0 || issue == 0 {

		url := content.extractByKey("url")
		if url == "" {
			return errors.New("invalid issue url")
		}
		owner, repo, issue, err = git.ParseIssueUrl(url)
	}
	if err != nil {
		return err
	}

	version := content.extractByKey("version")

	// TODO while the cmd is openned to R&D, we should check the user's permission
	// isForce := utils.Contains(content.flags, "--force")

	return service.ForceTriageIssue(version, owner, repo, issue, entity.VersionTriageResultAccept)
}
