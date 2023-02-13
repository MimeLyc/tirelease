package model

import (
	"tirelease/commons/git"
	"tirelease/internal/entity"
)

type PullRequest struct {
	*entity.PullRequest
	Author User `json:"author,omitempty"`
}

func NewPullRequest(entityPr entity.PullRequest) PullRequest {
	return PullRequest{
		PullRequest: &entityPr,
	}
}

func (pr PullRequest) Approve() error {
	if err := pr.Unlabel(git.NotCheryyPickLabel); err != nil {
		return err
	}

	if err := pr.Label(git.CherryPickLabel); err != nil {
		return err
	}

	return nil
}

func (pr PullRequest) UnApprove() error {
	if err := pr.Unlabel(git.CherryPickLabel); err != nil {
		return err
	}

	if err := pr.Label(git.NotCheryyPickLabel); err != nil {
		return err
	}

	return nil
}

func (pr PullRequest) Close() error {
	return git.ClientV4.ClosePullRequestsById(pr.PullRequestID)
}

func (pr PullRequest) Label(label string) error {
	// add issue label
	_, _, err := git.Client.AddLabel(pr.Owner, pr.Repo, pr.Number, label)
	if nil != err {
		return err
	}
	return nil
}

func (pr PullRequest) Unlabel(label string) error {
	// remove issue label
	_, err := git.Client.RemoveLabel(pr.Owner, pr.Repo, pr.Number, label)
	if nil != err {
		return err
	}
	return nil
}

// IsClosed method return true if pull request is closed and not merged
func (pr PullRequest) IsClosed() bool {
	return pr.State == "closed" && !pr.Merged
}

func (pr PullRequest) IsMerged() bool {
	return pr.Merged
}

// IsFinished method return true if pull request is closed or merged
func (pr PullRequest) IsFinished() bool {
	return pr.State == "closed" || pr.Merged
}
