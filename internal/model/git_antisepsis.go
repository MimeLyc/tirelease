package model

import (
	"sort"
	"time"
	"tirelease/commons/git"

	"tirelease/commons/utils"
)

type GitCommit struct {
	Oid            string
	AbbreviatedOid string
	CommittedTime  time.Time
	PushedTime     time.Time
}

type ByCommitTime []GitCommit

func (a ByCommitTime) Len() int           { return len(a) }
func (a ByCommitTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCommitTime) Less(i, j int) bool { return a[i].CommittedTime.After(a[j].CommittedTime) }

type GitRef struct {
	FirstCommit   GitCommit
	Owner         string
	Repo          string
	Name          string
	QualifiedName string
	PushedTime    time.Time
}

func GetUserByGitCode(clientId, clientSecret, code string) (*User, error) {
	accessToken, err := git.GetAccessTokenByClient(clientId, clientSecret, code)
	if err != nil {
		return nil, err
	}

	user, err := git.GetUserByToken(accessToken)
	if err != nil {
		return nil, err
	}

	employee, err := UserCmd{}.BuildByGhLogin(user.GetLogin())
	employee.GitUser = GitUser{
		GitID:        user.GetID(),
		GitLogin:     user.GetLogin(),
		GitAvatarURL: user.GetAvatarURL(),
		GitName:      user.GetName(),
	}

	return employee, nil
}

// Get checkout commit of ref(tag or branch).
// owner: org of the repo
// Warning: will trace back **501** commit from the ref endpoint.
// If there is no checkout commit found from the tracing back
//
//	the return result will be nil.
func GetCheckoutCommitOfRef(owner, repo, refName string, refType git.RefType) (*GitCommit, error) {
	commitFields, err := git.ClientV4.GetCommitsByQualifiedRef(owner, repo, refName, refType, nil, nil)
	if err != nil || len(commitFields) == 0 {
		return nil, err
	}
	refCommits := MapToGitCommit(commitFields)
	sort.Sort(ByCommitTime(refCommits))
	return getCheckoutCommit(owner, repo, refCommits), nil
}

// Get checkout commit in the range of refCommits.
// The refCommits should be sorted in commitDate order.
// Will trace back 501 commits from the **second** commit of refCommits.
//
//	so if the first commit of refCommits is in default branch,
//	the return result will be the second one.
func getCheckoutCommit(owner, repo string, refCommits []GitCommit) *GitCommit {
	startCommit := refCommits[0]

	lastDefaultCommitField, err := git.ClientV4.GetLastCommitOfDefaultBranchUntil(owner, repo, startCommit.CommittedTime)
	if err != nil {
		return nil
	}

	defaultCommits := MapToGitCommit([]git.CommitFiled{*lastDefaultCommitField})
	thresholdPage := 5
	pageCnt := 0

	for i := 0; i < len(defaultCommits) && pageCnt < thresholdPage; i++ {
		commit := defaultCommits[i]
		if utils.Contains(refCommits, commit) {
			return &commit
		}

		if i == len(defaultCommits)-1 {
			tmpFields, _, err := git.ClientV4.GetCommitsOfDefaultBranch(owner, repo, commit.Oid)
			if err != nil {
				return nil
			}
			defaultCommits = MapToGitCommit(tmpFields)
			sort.Sort(ByCommitTime(defaultCommits))
			i = -1
			pageCnt++
		}
	}

	return nil
}

func MapToGitCommit(commits []git.CommitFiled) []GitCommit {
	result := make([]GitCommit, 0)
	for _, commit := range commits {
		result = append(result, GitCommit{
			Oid:            string(commit.Oid),
			AbbreviatedOid: string(commit.AbbreviatedOid),
			CommittedTime:  commit.CommittedDate.Time,
			PushedTime:     commit.PushedDate.Time,
		})
	}

	return result
}
