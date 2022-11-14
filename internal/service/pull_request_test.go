package service

import (
	"fmt"
	"testing"

	"tirelease/commons/database"
	"tirelease/commons/git"
	"tirelease/commons/utils"
	"tirelease/internal/entity"
	"tirelease/internal/model"
	"tirelease/internal/repository"
	"tirelease/internal/service/component"

	"github.com/stretchr/testify/assert"
)

func TestGetPullRequestByNumberFromV3(t *testing.T) {
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)

	pr, err := GetPRByNumberFromV3(git.TestOwner, git.TestRepo, git.TestPullRequestId)

	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, pr != nil)
}

func TestGetPullRequestRefIssuesByRegexFromV4(t *testing.T) {
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)

	pr, err := git.ClientV4.GetPullRequestByID(git.TestPullRequestNodeID)
	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, pr != nil)

	issueNumbers, err := GetPullRequestRefIssuesByRegexFromV4(pr)
	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, len(issueNumbers) > 0)
}
func TestRegexReferenceNumbers(t *testing.T) {
	s := "close #1"
	issueNumbers, err := RegexReferenceNumbers(s)
	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, len(issueNumbers) == 1)

	s = "close #10, #100, #1000"
	issueNumbers, err = RegexReferenceNumbers(s)
	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, len(issueNumbers) == 3)
}

// ---- Scripts for Cleaning PRs----
// ---- Will close prs and post comments
const prCommentWontFix = `
This pull request is closed because the related issue is triaged "Won't Fix".
If it's still needed, you can reopen it or just regenerate it using bot,
see:
- https://prow.tidb.io/command-help#cherrypick
- https://book.prow.tidb.net/#/plugins/cherrypicker

You can find more details at:
- https://internals.tidb.io/t/topic/785
`

const prCommentEoSDVersions = `
This pull request is closed because it's related version has closed automatic cherry-picking.
If it's still needed, you can reopen it or just regenerate it using bot,
see:
- https://prow.tidb.io/command-help#cherrypick
- https://book.prow.tidb.net/#/plugins/cherrypicker

You can find more details at:
- https://internals.tidb.io/t/topic/785
`

func TestClosePrsOfTargetBranch(t *testing.T) {
	targetBranchs := []string{
		fmt.Sprintf("%s%s", git.ReleaseBranchPrefix, "3.0"),
		fmt.Sprintf("%s%s", git.ReleaseBranchPrefix, "4.0"),
		fmt.Sprintf("%s%s", git.ReleaseBranchPrefix, "5.0"),
		fmt.Sprintf("%s%s", git.ReleaseBranchPrefix, "5.1"),
		fmt.Sprintf("%s%s", git.ReleaseBranchPrefix, "5.2"),
		fmt.Sprintf("%s%s", git.ReleaseBranchPrefix, "6.0"),
		fmt.Sprintf("%s%s", git.ReleaseBranchPrefix, "6.2"),
	}

	// Init
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	database.Connect(generateConfig())

	prsToClose := make([]entity.PullRequest, 0)
	targetComponent := component.TIFLASH_COMPUTE

	prs := make([]entity.PullRequest, 0)
	for _, targetBranch := range targetBranchs {
		tmpPrs, _ := repository.SelectPullRequest(
			&entity.PullRequestOption{
				BaseBranch: targetBranch,
				State:      "open",
			},
		)
		prs = append(prs, *tmpPrs...)
	}

	for _, pr := range prs {
		relations, _ := repository.SelectIssuePrRelation(
			&entity.IssuePrRelationOption{
				PullRequestID: pr.PullRequestID,
			},
		)

		if len(*relations) == 0 {
			continue
		}

		allIssuesMatchComponent := true
		for _, relation := range *relations {
			issue, _ := repository.SelectIssueUnique(
				&entity.IssueOption{
					IssueID: relation.IssueID,
				},
			)
			if !utils.Contains(component.GetComponents(issue.Owner, issue.Repo, issue.LabelsString), targetComponent) {
				allIssuesMatchComponent = false
			}
		}
		if allIssuesMatchComponent {
			prsToClose = append(prsToClose, pr)
		}
	}

	for _, pr := range prsToClose {
		ClosePRWithComment(pr.PullRequestID, pr.Owner, pr.Repo, pr.Number, prCommentEoSDVersions)
	}
}

// script to close won't fix pr of target component
func TestClosePrsOfWontFix(t *testing.T) {

	// Init
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	database.Connect(generateConfig())

	prsToClose := make([]entity.PullRequest, 0)
	// targetComponent := component.TIDB_PLANNER

	versionTriages, err := repository.SelectVersionTriage(
		&entity.VersionTriageOption{
			TriageResult: entity.VersionTriageResultWontFix,
		},
	)
	assert.Nil(t, err)

	for _, triage := range *versionTriages {
		issueID := triage.IssueID
		versionName := triage.VersionName
		minorVersionName := model.ExtractVersionMinorName(versionName)
		// issue, _ := repository.SelectIssueUnique(
		// 	&entity.IssueOption{
		// 		IssueID: issueID,
		// 	},
		// )
		// if !utils.Contains(component.GetComponents(issue.Owner, issue.Repo, issue.LabelsString), targetComponent) {
		// 	continue
		// }

		baseBranch := fmt.Sprintf("%s%s", git.ReleaseBranchPrefix, minorVersionName)

		relatedPrs, err := repository.SelectIssuePrRelation(
			&entity.IssuePrRelationOption{
				IssueID: issueID,
			},
		)
		assert.Nil(t, err)
		if len(*relatedPrs) == 0 {
			continue
		}

		prIds := make([]string, 0)

		for _, relation := range *relatedPrs {
			prIds = append(prIds, relation.PullRequestID)
		}
		prs, err := repository.SelectPullRequest(
			&entity.PullRequestOption{
				PullRequestIDs: prIds,
				State:          "open",
				BaseBranch:     baseBranch,
			},
		)
		assert.Nil(t, err)
		for _, pr := range *prs {
			prsToClose = append(prsToClose, pr)
		}

	}

	for _, pr := range prsToClose {
		ClosePRWithComment(pr.PullRequestID, pr.Owner, pr.Repo, pr.Number, prCommentWontFix)
	}
}

func ClosePRWithComment(prId, owner, repo string, number int, comment string) {
	if err := git.ClientV4.ClosePullRequestsById(prId); err != nil {
		fmt.Printf("%v", err)
	}
	_, _, err := git.Client.CreateCommentByNumber(owner, repo, number, comment)
	if err != nil {
		panic(err)
	}
}
