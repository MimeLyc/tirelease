package service

import (
	"fmt"
	"tirelease/internal/entity"
	"tirelease/internal/model"
)

func NotifySprintIssueMetrics(request entity.SprintMetaRequest, user model.User) error {
	masterIssues, err := model.SelectIssuesFixedBeforeSprintCheckout(*request.Major, *request.Minor)
	// Get all Issues fixed before sprint checkout.
	if err != nil {
		return err
	}

	branchIssues, err := model.SelectIssuesFixedAfterSprintCheckout(
		*request.Major,
		*request.Minor,
		entity.IssueOption{
			State: "closed",
		},
	)
	if err != nil {
		return err
	}

	fmt.Sprintf("%v", masterIssues)
	fmt.Sprintf("%v", branchIssues)

	// query issues info
	//   fixed before frozen, fixed after frozen
	// dump file and upload
	// make norification
	return nil
}
