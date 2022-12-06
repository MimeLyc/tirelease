package service

import (
	"fmt"
	"time"
	"tirelease/commons/fileserver"
	"tirelease/commons/ifile"
	"tirelease/internal/model"
	"tirelease/internal/service/notify"
)

func NotifyPatchReleaseNotesExcel(major, minor, patch int, email string) error {
	prs, err := model.SelectMergedPrsOfSprint(major, minor)
	if err != nil {
		return err
	}

	prIssueRelation, err := model.PrIssueRelationBuilder{}.BuildByPrs(prs)

	if err != nil {
		return err
	}

	releaseNotePrs := model.DumpReleaseNotePullRequests(prIssueRelation)

	sprintName := ComposeVersionMinorNameByNumber(major, minor)
	timestamp := time.Now().Format("2006-01-02")
	sprint_file_prefix := fmt.Sprintf(SprintReleaseNoteExcelFilePrefix, sprintName)
	filename := fmt.Sprintf(TmpFileFormat, sprint_file_prefix, timestamp, ExcelPostFix)
	qualifiedName := fmt.Sprintf("%s/%s", TmpFileDir, filename)
	defer ifile.RmAllFile(TmpFileDir)

	err = ifile.CreateExcelSheetByTag(releaseNotePrs, TmpFileDir, filename, SprintReleaseNoteExcelSheet)
	if err != nil {
		return err
	}
	downloadUrl, err := fileserver.UploadFile(qualifiedName, fmt.Sprintf("%s/%s", TiReleaseFileServerTmpDir, filename))

	if err != nil {
		return err
	}

	content := composeSprintReleaseNoteNotifyContent(sprintName, filename, downloadUrl)
	err = notify.SendFeishuFormattedByEmail(email, content)
	return err
}

func NotifySprintReleaseNotesExcel(major, minor int, email string) error {
	prs, err := model.SelectMergedPrsOfSprint(major, minor)
	if err != nil {
		return err
	}

	prIssueRelation, err := model.PrIssueRelationBuilder{}.BuildByPrs(prs)

	if err != nil {
		return err
	}

	releaseNotePrs := model.DumpReleaseNotePullRequests(prIssueRelation)

	sprintName := ComposeVersionMinorNameByNumber(major, minor)
	timestamp := time.Now().Format("2006-01-02")
	sprint_file_prefix := fmt.Sprintf(SprintReleaseNoteExcelFilePrefix, sprintName)
	filename := fmt.Sprintf(TmpFileFormat, sprint_file_prefix, timestamp, ExcelPostFix)
	qualifiedName := fmt.Sprintf("%s/%s", TmpFileDir, filename)
	defer ifile.RmAllFile(TmpFileDir)

	err = ifile.CreateExcelSheetByTag(releaseNotePrs, TmpFileDir, filename, SprintReleaseNoteExcelSheet)
	if err != nil {
		return err
	}
	downloadUrl, err := fileserver.UploadFile(qualifiedName, fmt.Sprintf("%s/%s", TiReleaseFileServerTmpDir, filename))

	if err != nil {
		return err
	}

	content := composeSprintReleaseNoteNotifyContent(sprintName, filename, downloadUrl)
	err = notify.SendFeishuFormattedByEmail(email, content)
	return err
}

const sprintReleaseNoteNorifyContentText = `
The following link is the prs merged in sprint %s with formated strings for release note.
Granularity of each row in the excel file is **pr + issue**.
Please click to download:
`

func composeSprintReleaseNoteNotifyContent(sprintname, filename, downloadUrl string) notify.NotifyContent {
	link := notify.Link{
		Href: downloadUrl,
		Text: filename,
	}
	block := notify.Block{
		Text:  fmt.Sprintf(sprintReleaseNoteNorifyContentText, sprintname),
		Links: []notify.Link{link},
	}
	content := notify.NotifyContent{
		Header: fmt.Sprintf("Pullrequest Dumped For Release Note in Sprint %s", sprintname),
		Blocks: []notify.Block{block},
	}

	return content
}
