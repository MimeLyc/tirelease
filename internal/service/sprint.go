package service

import (
	"fmt"
	"time"
	"tirelease/commons/fileserver"
	"tirelease/commons/git"
	"tirelease/commons/ifile"
	"tirelease/internal/dto"
	"tirelease/internal/entity"
	"tirelease/internal/model"
	"tirelease/internal/service/notify"
)

func FindSprintIssues(major, minor int, option entity.IssueOption) (*dto.SprintIssuesResponse, error) {
	masterIssues, err := model.SelectIssuesBeforeSprintCheckout(major, minor, option)
	// Get all Issues fixed before sprint checkout.
	if err != nil {
		return nil, err
	}

	masterIssueDtos := make([]dto.SprintIssue, 0)
	for _, issue := range masterIssues {
		issue := issue
		masterIssueDtos = append(masterIssueDtos,
			dto.SprintIssue{
				Issue:   issue,
				IsBlock: false,
			},
		)
	}

	branchIssues, err := model.SelectIssuesAfterSprintCheckout(
		major,
		minor,
		option,
	)
	if err != nil {
		return nil, err
	}

	branchIssueDtos := make([]dto.SprintIssue, 0)
	for _, issue := range branchIssues {
		issue := issue
		branchIssueDtos = append(branchIssueDtos,
			dto.SprintIssue{
				Issue: issue,
				IsBlock: model.GetBlockDefaultValue(
					&issue.Issue,
					ComposeVersionMinorNameByNumber(major, minor),
					&issue.VersionTriages,
				) == entity.BlockVersionReleaseResultBlock,
			},
		)
	}

	return &dto.SprintIssuesResponse{
		Major:        major,
		Minor:        minor,
		MasterIssues: &masterIssueDtos,
		BranchIssues: &branchIssueDtos,
	}, nil
}

func NotifySprintBugMetrics(major, minor int, email string) error {
	issueOption := entity.IssueOption{
		TypeLabel: git.BugTypeLabel,
	}

	masterIssues, err := model.SelectIssuesBeforeSprintCheckout(major, minor, issueOption)
	// Get all Issues fixed before sprint checkout.
	if err != nil {
		return err
	}

	branchIssues, err := model.SelectIssuesAfterSprintCheckout(
		major,
		minor,
		issueOption,
	)
	if err != nil {
		return err
	}

	sprintName := ComposeVersionMinorNameByNumber(major, minor)
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	sprint_file_prefix := fmt.Sprintf(SprintFixedIssueMetricFilePrefix, sprintName)
	filename := fmt.Sprintf(TmpFileFormat, sprint_file_prefix, timestamp, ExcelPostFix)
	qualifiedName := fmt.Sprintf("%s/%s", TmpFileDir, filename)
	defer ifile.RmAllFile(TmpFileDir)
	err = ifile.CreateExcelSheetByTag(masterIssues, TmpFileDir, filename, SprintFixedIssueMetricMasterSheet)
	if err != nil {
		return err
	}

	if len(branchIssues) > 0 {
		err = ifile.CreateExcelSheetByTag(branchIssues, TmpFileDir, filename, SprintFixedIssueMetricBranchSheet)
		if err != nil {
			return err
		}
	}

	downloadUrl, err := fileserver.UploadFile(qualifiedName, fmt.Sprintf("%s/%s", TiReleaseFileServerTmpDir, filename))

	if err != nil {
		return err
	}

	content := composeSprintBugMetricNotifyContent(sprintName, filename, downloadUrl)
	err = notify.SendFeishuFormattedByEmail(email, content)
	return err
}

const sprintBugMetricNorifyContentText = `
The following link is the bug issues fixed before sprint frozen( in the "master_fixed" sheet ) and after sprint frozen( in the "branch_fixed" sheet ), please click to download:
`

func composeSprintBugMetricNotifyContent(sprintname, filename, downloadUrl string) notify.NotifyContent {
	link := notify.Link{
		Href: downloadUrl,
		Text: filename,
	}
	block := notify.Block{
		Text:  sprintBugMetricNorifyContentText,
		Links: []notify.Link{link},
	}
	content := notify.NotifyContent{
		Header: fmt.Sprintf("Bug fixed Details of Sprint %s", sprintname),
		Blocks: []notify.Block{block},
	}

	return content
}
