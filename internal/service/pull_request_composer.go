package service

import (
	"strings"
	"time"
	"tirelease/commons/git"
	"tirelease/internal/entity"

	"github.com/google/go-github/v41/github"
)

func ComposePullRequest(mapper GitPullRequestMapper) *entity.PullRequest {
	pr := &entity.PullRequest{
		PullRequestID:      mapper.NodeID(),
		Number:             mapper.Number(),
		Title:              mapper.Title(),
		State:              mapper.State(),
		Owner:              mapper.Owner(),
		Repo:               mapper.Repo(),
		BaseBranch:         mapper.BaseBranch(),
		HTMLURL:            mapper.HTMLURL(),
		CreateTime:         mapper.CreateTime(),
		UpdateTime:         mapper.UpdateTime(),
		CloseTime:          mapper.CloseTime(),
		MergeTime:          mapper.MergeTime(),
		Merged:             mapper.Merged(),
		MergeableState:     mapper.MergeableState(),
		CherryPickApproved: mapper.CherryPickApproved(),
		AlreadyReviewed:    mapper.AlreadyReviewed(),
		Labels:             mapper.Labels(),
		Assignees:          mapper.Assignees(),
		RequestedReviewers: mapper.RequestedReviewers(),
		Body:               mapper.Body(),
	}

	releaseNote, _ := parseReleaseNote(*pr)
	pr.IsReleaseNoteConfirmed = releaseNote.IsReleaseNoteConfirmed
	pr.ReleaseNote = releaseNote.ReleaseNote

	return pr
}

func parseReleaseNote(prEntity entity.PullRequest) (git.ReleaseNoteData, error) {
	releaseNote, err := git.ParseReleaseNote(prEntity.Body)

	if err != nil || !releaseNote.IsReleaseNoteConfirmed {
		labels := prEntity.Labels
		for _, label := range *labels {
			if *label.Name == git.NoneReleaseNoteLabel {
				releaseNote.IsReleaseNoteConfirmed = true
				releaseNote.ReleaseNote = "None"
			}
		}
	}

	return releaseNote, err
}

func ComposePRFromV3(pr *github.PullRequest) *entity.PullRequest {
	mapper := &GitPullRequestMapperV3{
		PullRequest: pr,
	}
	return ComposePullRequest(mapper)
}

// Query PullRequest From Github And Construct Issue Data Service
func GetPRByNumberFromV3(owner, repo string, number int) (*entity.PullRequest, error) {
	pr, _, err := git.Client.GetPullRequestByNumber(owner, repo, number)
	if nil != err {
		return nil, err
	}
	return ComposePRFromV3(pr), nil
}

func GetPRsByRequestFromV4(request *git.RemoteIssueRangeRequest) ([]entity.PullRequest, error) {
	prFileds, err := git.ClientV4.GetPullRequestsFromV4(request)
	if err != nil {
		return nil, err
	}
	prs := make([]entity.PullRequest, 0)

	for _, field := range prFileds {
		mapper := &GitPullRequestMapperV4{
			PullRequest: &field,
		}
		prs = append(prs, *ComposePullRequest(mapper))
	}
	return prs, nil
}

func ComposePullRequestWithoutTimelineFromV4(withoutTimeline *git.PullRequestFieldWithoutTimelineItems) *entity.PullRequest {
	pullRequestField := &git.PullRequestField{
		PullRequestFieldWithoutTimelineItems: *withoutTimeline,
	}
	mapper := &GitPullRequestMapperV4{
		PullRequest: pullRequestField,
	}
	return ComposePullRequest(mapper)
}

type GitPullRequestMapper interface {
	NodeID() string
	Number() int
	Title() string
	State() string
	Owner() string
	Repo() string
	BaseBranch() string
	HTMLURL() string
	CreateTime() time.Time
	UpdateTime() time.Time
	CloseTime() *time.Time
	MergeTime() *time.Time
	Merged() bool
	MergeableState() *string
	CherryPickApproved() bool
	AlreadyReviewed() bool
	Labels() *[]github.Label
	Assignees() *[]github.User
	RequestedReviewers() *[]github.User
	Body() string
}

type GitPullRequestMapperV4 struct {
	PullRequest *git.PullRequestField
}

func (m *GitPullRequestMapperV4) NodeID() string {
	return m.PullRequest.ID.(string)
}

func (m *GitPullRequestMapperV4) Number() int {
	return int(m.PullRequest.Number)
}

func (m *GitPullRequestMapperV4) Title() string {
	return string(m.PullRequest.Title)
}

func (m *GitPullRequestMapperV4) State() string {
	return strings.ToLower(string(m.PullRequest.State))
}

func (m *GitPullRequestMapperV4) Owner() string {
	return string(m.PullRequest.Repository.Owner.Login)
}

func (m *GitPullRequestMapperV4) Repo() string {
	return string(m.PullRequest.Repository.Name)
}

func (m *GitPullRequestMapperV4) BaseBranch() string {
	return string(m.PullRequest.BaseRefName)
}

func (m *GitPullRequestMapperV4) HTMLURL() string {
	return string(m.PullRequest.Url)
}

func (m *GitPullRequestMapperV4) CreateTime() time.Time {
	return m.PullRequest.CreatedAt.Time
}

func (m *GitPullRequestMapperV4) UpdateTime() time.Time {
	return m.PullRequest.UpdatedAt.Time
}

func (m *GitPullRequestMapperV4) CloseTime() *time.Time {
	closeAt := m.PullRequest.ClosedAt
	if closeAt == nil {
		return nil
	}
	return &closeAt.Time
}

func (m *GitPullRequestMapperV4) MergeTime() *time.Time {
	mergeAt := m.PullRequest.MergedAt
	if mergeAt == nil {
		return nil
	}
	return &mergeAt.Time
}

func (m *GitPullRequestMapperV4) Merged() bool {
	return bool(m.PullRequest.Merged)
}

func (m *GitPullRequestMapperV4) MergeableState() *string {
	mergeableState := strings.ToLower(string(m.PullRequest.Mergeable))
	return &mergeableState
}

func (m *GitPullRequestMapperV4) CherryPickApproved() bool {
	cherryPickApproved := false
	for i := range m.PullRequest.Labels.Nodes {
		node := m.PullRequest.Labels.Nodes[i]
		label := github.Label{
			Name: github.String(string(node.Name)),
		}

		if *label.Name == git.CherryPickLabel {
			cherryPickApproved = true
		}
	}

	return cherryPickApproved
}

func (m *GitPullRequestMapperV4) AlreadyReviewed() bool {
	alreadyReviwed := false
	for i := range m.PullRequest.Labels.Nodes {
		node := m.PullRequest.Labels.Nodes[i]
		label := github.Label{
			Name: github.String(string(node.Name)),
		}
		if *label.Name == git.LGT2Label {
			alreadyReviwed = true
		}
	}
	return alreadyReviwed
}

func (m *GitPullRequestMapperV4) Labels() *[]github.Label {
	labels := &[]github.Label{}
	for i := range m.PullRequest.Labels.Nodes {
		node := m.PullRequest.Labels.Nodes[i]
		label := github.Label{
			Name: github.String(string(node.Name)),
		}
		*labels = append(*labels, label)

	}
	return labels
}

func (m *GitPullRequestMapperV4) Assignees() *[]github.User {
	assignees := &[]github.User{}
	for i := range m.PullRequest.Assignees.Nodes {
		node := m.PullRequest.Assignees.Nodes[i]
		user := github.User{
			Login: (*string)(&node.Login),
		}
		*assignees = append(*assignees, user)
	}

	return assignees
}

func (m *GitPullRequestMapperV4) RequestedReviewers() *[]github.User {
	requestedReviewers := &[]github.User{}
	for i := range m.PullRequest.ReviewRequests.Nodes {
		node := m.PullRequest.ReviewRequests.Nodes[i]
		user := github.User{
			Login: (*string)(&node.RequestedReviewer.Login),
		}
		*requestedReviewers = append(*requestedReviewers, user)
	}
	return requestedReviewers
}

func (m *GitPullRequestMapperV4) Body() string {
	return string(m.PullRequest.Body)
}

type GitPullRequestMapperV3 struct {
	PullRequest *github.PullRequest
}

func (m *GitPullRequestMapperV3) NodeID() string {
	return *m.PullRequest.NodeID
}

func (m *GitPullRequestMapperV3) Number() int {
	return *m.PullRequest.Number
}

func (m *GitPullRequestMapperV3) Title() string {
	return *m.PullRequest.Title
}

func (m *GitPullRequestMapperV3) State() string {
	return *m.PullRequest.State
}

func (m *GitPullRequestMapperV3) Owner() string {
	return *m.PullRequest.Base.Repo.Owner.Login
}

func (m *GitPullRequestMapperV3) Repo() string {
	return *m.PullRequest.Base.Repo.Name
}

func (m *GitPullRequestMapperV3) BaseBranch() string {
	return *m.PullRequest.Base.Ref
}

func (m *GitPullRequestMapperV3) HTMLURL() string {
	return *m.PullRequest.HTMLURL
}

func (m *GitPullRequestMapperV3) CreateTime() time.Time {
	return *m.PullRequest.CreatedAt
}

func (m *GitPullRequestMapperV3) UpdateTime() time.Time {
	return *m.PullRequest.UpdatedAt
}

func (m *GitPullRequestMapperV3) CloseTime() *time.Time {
	return m.PullRequest.ClosedAt
}

func (m *GitPullRequestMapperV3) MergeTime() *time.Time {
	return m.PullRequest.MergedAt
}

func (m *GitPullRequestMapperV3) Merged() bool {
	return *m.PullRequest.Merged
}

func (m *GitPullRequestMapperV3) MergeableState() *string {
	mergeState := strings.ToLower(*m.PullRequest.MergeableState)
	return &mergeState
}

func (m *GitPullRequestMapperV3) CherryPickApproved() bool {
	cherryPickApproved := false
	for i := range m.PullRequest.Labels {
		node := m.PullRequest.Labels[i]
		label := github.Label{
			Name:  node.Name,
			Color: node.Color,
		}

		if *label.Name == git.CherryPickLabel {
			cherryPickApproved = true
		}
	}
	return cherryPickApproved
}

func (m *GitPullRequestMapperV3) AlreadyReviewed() bool {
	alreadyReviwed := false
	for i := range m.PullRequest.Labels {
		node := m.PullRequest.Labels[i]
		label := github.Label{
			Name:  node.Name,
			Color: node.Color,
		}

		if *label.Name == git.LGT2Label {
			alreadyReviwed = true
		}
	}
	return alreadyReviwed
}

func (m *GitPullRequestMapperV3) Labels() *[]github.Label {
	labels := &[]github.Label{}
	for i := range m.PullRequest.Labels {
		node := m.PullRequest.Labels[i]
		label := github.Label{
			Name:  node.Name,
			Color: node.Color,
		}
		*labels = append(*labels, label)
	}
	return labels
}

func (m *GitPullRequestMapperV3) Assignees() *[]github.User {
	assignees := &[]github.User{}
	for i := range m.PullRequest.Assignees {
		node := m.PullRequest.Assignees[i]
		user := github.User{
			Login: node.Login,
		}
		*assignees = append(*assignees, user)
	}
	return assignees
}

func (m *GitPullRequestMapperV3) RequestedReviewers() *[]github.User {
	requestedReviewers := &[]github.User{}
	for i := range m.PullRequest.RequestedReviewers {
		node := m.PullRequest.RequestedReviewers[i]
		user := github.User{
			Login: node.Login,
		}
		*requestedReviewers = append(*requestedReviewers, user)
	}
	return requestedReviewers
}

func (m *GitPullRequestMapperV3) Body() string {
	if m.PullRequest.Body == nil {
		return ""
	}
	return *m.PullRequest.Body
}
