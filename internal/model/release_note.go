package model

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

	Owner                    string `excel:"repo_org"`
	Repo                     string `excel:"repo"`
	IsPrRefMultiIssue        bool   `excel:"is_pr_ref_multi_issue"`
	PrAuthorName             string `excel:"pr_author_name"`
	IsPrAuthorActiveEmployee bool   `excel:"is_pr_author_active_employee"`
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

func NewReleaseNotePullRequest(pr PullRequest, issue *Issue) ReleaseNotePullRequest {
	components := ComponentString(pr, issue)
	formatedReleaseNote := FormatedReleaseNote(pr, issue)

	result := ReleaseNotePullRequest{
		Owner:                    pr.Owner,
		Repo:                     pr.Repo,
		Components:               components,
		State:                    pr.State,
		Title:                    pr.Title,
		HTMLURL:                  pr.HTMLURL,
		ReleaseNote:              pr.ReleaseNote,
		FormatedReleaseNote:      formatedReleaseNote,
		PrAuthor:                 pr.AuthorGhLogin,
		PrAuthorName:             pr.Author.Name,
		IsPrAuthorActiveEmployee: pr.Author.IsActive && pr.Author.IsEmployee,
	}

	if issue == nil {
		return result
	}

	formatedReleaseNote = FormatedReleaseNote(pr, issue)
	result.IssueSeverity = issue.SeverityLabel
	result.IssueType = issue.TypeLabel
	result.FormatedReleaseNote = formatedReleaseNote

	return result
}
