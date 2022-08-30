package model

import (
	"tirelease/commons/git"
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

// Hook github api to change the `approve` related labels.
func ChangePrApprovedLabel(pr entity.PullRequest, isFrozen, isAccept bool) error {
	if !isFrozen && isAccept {
		RemoveLabelByPullRequestID(pr, git.NotCheryyPickLabel)

		err := AddLabelByPullRequestID(pr, git.CherryPickLabel)
		if err != nil {
			return err
		}
	} else {
		RemoveLabelByPullRequestID(pr, git.CherryPickLabel)

		err := AddLabelByPullRequestID(pr, git.NotCheryyPickLabel)
		if err != nil {
			return err
		}
	}

	return nil
}

func RemoveLabelByPullRequestID(pr entity.PullRequest, label string) error {
	// remove issue label
	_, err := git.Client.RemoveLabel(pr.Owner, pr.Repo, pr.Number, label)
	if nil != err {
		return err
	}
	return nil
}

// Operation
func AddLabelByPullRequestID(pr entity.PullRequest, label string) error {
	// add issue label
	_, _, err := git.Client.AddLabel(pr.Owner, pr.Repo, pr.Number, label)
	if nil != err {
		return err
	}
	return nil
}

func closeWontfixPrs(prs []entity.PullRequest, originIssueID, versionName string) error {
	for _, pr := range prs {
		isAllWontFix := true
		relations, err := repository.SelectIssuePrRelation(
			&entity.IssuePrRelationOption{
				PullRequestID: pr.PullRequestID,
			},
		)
		if err != nil {
			return err
		}

		// Close pr only when all related issues were closed in the version
		for _, relation := range *relations {
			issueId := relation.IssueID
			if issueId == originIssueID {
				continue
			}
			triage, err := repository.SelectVersionTriageUnique(
				&entity.VersionTriageOption{
					IssueID: issueId,
					// TODO: check all patch versions under the same minor version.
					// The reason why not do it now(20220830) is that it'll hard to find all related triage history which is not in the same patch version.
					VersionName: versionName,
				},
			)
			if err != nil {
				return err
			}
			if triage != nil && triage.TriageResult != entity.VersionTriageResultWontFix {
				isAllWontFix = false
			}
		}

		if isAllWontFix {
			git.ClientV4.ClosePullRequestsById(pr.PullRequestID)
		}
	}

	return nil
}
