package model

import (
	"tirelease/internal/entity"
)

// 1:1 PullRequest: Issue
type ReleaseNotePullRequest struct {
	// Fileds bellow is pointed by Doc Team.
	HTMLURL             string `excel:"pr_link"`
	Components          string `excel:"components"`
	Title               string `excel:"pr_title"`
	State               string `excel:"pr_status"`
	PrAuthor            string `excel:"pr_author"`
	ReleaseNote         string `excel:"release_note"`
	FormatedReleaseNote string `excel:"formated_release_note"`
	IssueSeverity       string `excel:"issue_severity"`
	IssueType           string `excel:"issue_type"`

	Owner             string `excel:"repo_org"`
	Repo              string `excel:"repo"`
	IsPrRefMultiIssue bool   `excel:"is_pr_ref_multi_issue"`

	// Commented fileds bellow is redundant for further using untile the final fields for release notes is confirmed.

	// PullRequestID string `excel:"pr_id"`
	// PrNumber      int    `excel:"pr_number"`
	// BaseBranch    string `excel:"pr_base_branch"`
	//
	// CreateTime time.Time  `excel:"pr_create_time"`
	// CloseTime  *time.Time `excel:"pr_close_time"`
	// MergeTime  *time.Time `excel:"pr_merge_time"`
	//
	// Merged bool `excel:"is_pr_merged"`
	//
	// PrLabels               string `excel:"pr_labels_string"`
	// PrAssignees            string `excel:"assignees_github_ids"`
	// IsReleaseNoteConfirmed bool   `excel:"has_relate_note"`
	//
	// // Issue columns
	// IssueID        string `excel:"issue_id"`
	// IssueLabel     string `excel:"issue_labels"`
	// IssueUrl       string `excel:"issue_url"`
	// IssueAssignees string `excel:"issue_assignees"`
	//
	// // release note coalesce([#issue_number](issue_url), [#pr_number](pr_url)) @assignee
	// FormatedPrAuthor       string `excel:"formated_pr_author"`
	// FormatedPrAssignees    string `excel:"formated_pr_assignees"`
	// FormatedIssueAssignees string `excel:"formated_issue_assignees"`
	// FormatedPrUrl          string `excel:"formated_pr_url"`
	// FormatedIssueUrl       string `excel:"formated_issue_url"`
}

func NewReleaseNotePullRequest(pr PullRequest, issue *entity.Issue) ReleaseNotePullRequest {
	components := ComponentString(pr, issue)
	// prLabels := PullRequestLabelsString(pr)
	// prAssignees := PullRequestAssigneesString(pr)
	// prFormatedAuthor := FormatedPullRequestAuthorString(pr)
	// prFormatedAssignees := FormatedPullRequestAssigneesString(pr)
	// prFormatedUrl := FormatedPullrequestUrl(pr)
	formatedReleaseNote := FormatedReleaseNote(pr, issue)

	result := ReleaseNotePullRequest{
		Owner:               pr.Owner,
		Repo:                pr.Repo,
		Components:          components,
		State:               pr.State,
		Title:               pr.Title,
		HTMLURL:             pr.HTMLURL,
		ReleaseNote:         pr.ReleaseNote,
		FormatedReleaseNote: formatedReleaseNote,
		PrAuthor:            pr.AuthorGhLogin,
		// PullRequestID:          pr.PullRequestID,
		// PrNumber:               pr.Number,
		// BaseBranch:             pr.BaseBranch,
		// CreateTime:             pr.CreateTime,
		// CloseTime:              pr.CloseTime,
		// MergeTime:              pr.MergeTime,
		// Merged:                 pr.Merged,
		// PrLabels:               prLabels,
		// PrAssignees:            prAssignees,
		// IsReleaseNoteConfirmed: pr.IsReleaseNoteConfirmed,
		// FormatedPrAuthor:       prFormatedAuthor,
	}

	if issue == nil {
		return result
	}

	formatedReleaseNote = FormatedReleaseNote(pr, issue)
	result.IssueSeverity = issue.SeverityLabel
	result.IssueType = issue.TypeLabel
	// issueAssignees := IssueAssigneesString(*issue)
	// issueLabels := IssueLabelsString(*issue)
	// issueFormatedAssignees := FormatedIssueAssigneesString(*issue)
	// issueFormatedUrl := FormatedIssueUrl(*issue)
	//
	// result.IssueID = issue.IssueID
	// result.IssueLabel = issueLabels
	// result.IssueUrl = issue.HTMLURL
	// result.IssueAssignees = issueAssignees
	// result.FormatedPrAssignees = prFormatedAssignees
	// result.FormatedIssueAssignees = issueFormatedAssignees
	// result.FormatedPrUrl = prFormatedUrl
	// result.FormatedIssueUrl = issueFormatedUrl
	result.FormatedReleaseNote = formatedReleaseNote

	return result
}

func DumpReleaseNotePullRequests(relations []PrIssueRelation) []ReleaseNotePullRequest {
	result := make([]ReleaseNotePullRequest, 0)
	for _, relation := range relations {
		pr := relation.PullRequest

		isPrRefMultiIssues := false
		if len(relation.RelatedIssues) == 0 {
			result = append(result, NewReleaseNotePullRequest(*pr, nil))
			continue
		} else if len(relation.RelatedIssues) > 1 {

			isPrRefMultiIssues = true
		}

		for _, issue := range relation.RelatedIssues {
			releaseNotePr := NewReleaseNotePullRequest(*pr, &issue)
			releaseNotePr.IsPrRefMultiIssue = isPrRefMultiIssues
			result = append(result, releaseNotePr)
		}

	}
	return result

}
