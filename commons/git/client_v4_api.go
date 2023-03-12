package git

import (
	"context"
	"fmt"
	"time"

	"github.com/shurcooL/githubv4"
)

// ============================================================================ Issue
func (client *GithubInfoV4) GetIssueByID(id string) (*IssueField, error) {
	var query struct {
		Node struct {
			IssueField `graphql:"... on Issue"`
		} `graphql:"node(id: $id)"`
	}
	params := map[string]interface{}{
		"id": githubv4.ID(id),
	}
	if err := client.client.Query(context.Background(), &query, params); err != nil {
		return nil, err
	}
	return &query.Node.IssueField, nil
}

func (client *GithubInfoV4) GetIssueWithoutTimelineByID(id string) (*IssueFieldWithoutTimelineItems, error) {
	var query struct {
		Node struct {
			IssueFieldWithoutTimelineItems `graphql:"... on Issue"`
		} `graphql:"node(id: $id)"`
	}
	params := map[string]interface{}{
		"id": githubv4.ID(id),
	}
	if err := client.client.Query(context.Background(), &query, params); err != nil {
		return nil, err
	}
	return &query.Node.IssueFieldWithoutTimelineItems, nil
}

func (client *GithubInfoV4) GetIssueByNumber(owner, name string, number int) (*IssueField, error) {
	var query struct {
		Repository struct {
			Issue struct {
				IssueField
			} `graphql:"issue(number: $number)"`
		} `graphql:"repository(name: $name, owner: $owner)"`
	}
	num := githubv4.Int(number)
	params := map[string]interface{}{
		"number": num,
		"name":   githubv4.String(name),
		"owner":  githubv4.String(owner),
	}
	if err := client.client.Query(context.Background(), &query, params); err != nil {
		return nil, err
	}
	return &query.Repository.Issue.IssueField, nil
}

// ============================================================================ PullRequest
func (client *GithubInfoV4) GetPullRequestByID(id string) (*PullRequestField, error) {
	var query struct {
		Node struct {
			PullRequestField `graphql:"... on PullRequest"`
		} `graphql:"node(id: $id)"`
	}
	params := map[string]interface{}{
		"id": githubv4.ID(id),
	}
	if err := client.client.Query(context.Background(), &query, params); err != nil {
		return nil, err
	}
	return &query.Node.PullRequestField, nil
}

func (client *GithubInfoV4) GetPullRequestWithoutTimelineByID(id string) (*PullRequestFieldWithoutTimelineItems, error) {
	var query struct {
		Node struct {
			PullRequestFieldWithoutTimelineItems `graphql:"... on PullRequest"`
		} `graphql:"node(id: $id)"`
	}
	params := map[string]interface{}{
		"id": githubv4.ID(id),
	}
	if err := client.client.Query(context.Background(), &query, params); err != nil {
		return nil, err
	}
	return &query.Node.PullRequestFieldWithoutTimelineItems, nil
}

func (client *GithubInfoV4) GetPullRequestsByNumber(owner, name string, number int) (*PullRequestField, error) {
	var query struct {
		Repository struct {
			PullRequest struct {
				PullRequestField
			} `graphql:"pullRequest(number: $number)"`
		} `graphql:"repository(name: $name, owner: $owner)"`
	}
	num := githubv4.Int(number)
	params := map[string]interface{}{
		"number": num,
		"name":   githubv4.String(name),
		"owner":  githubv4.String(owner),
	}
	if err := client.client.Query(context.Background(), &query, params); err != nil {
		return nil, err
	}
	return &query.Repository.PullRequest.PullRequestField, nil
}

func (client *GithubInfoV4) GetCommitsByQualifiedRef(owner, repo,
	name string, refType RefType, since, until *time.Time) ([]CommitFiled, error) {
	if refType == RefTypeTag {
		return client.GetHistoryCommitsByTag(owner, repo, name, since, until)
	} else if refType == RefTypeBranch {
		return client.GetCommitsByBranch(owner, repo, name, since, until)
	} else {
		return nil, fmt.Errorf("Error ref type of %s", refType)
	}
}

func (client *GithubInfoV4) GetCommitsByBranch(owner, repo, name string, since, until *time.Time) ([]CommitFiled, error) {
	name = fmt.Sprintf("refs/heads/%s", name)
	if since != nil || until != nil {
		return client.GetCommitsByRefWithTimeScope(owner, repo, name, since, until)
	} else {
		return client.GetCommitsByRef(owner, repo, name)
	}
}

func (client *GithubInfoV4) GetHistoryCommitsByTag(owner, repo, name string, since, until *time.Time) ([]CommitFiled, error) {
	name = fmt.Sprintf("refs/tags/%s", name)
	if since != nil || until != nil {
		return client.GetCommitsByRefWithTimeScope(owner, repo, name, since, until)
	} else {
		return client.GetCommitsByRef(owner, repo, name)
	}
}

// GetCommits history by ref
// History may trace back up to 64 commits due to the limitation of git graphql api
// Ref may be branch or tag or other git refs.
func (client *GithubInfoV4) GetCommitsByRefWithTimeScope(owner, repo, ref string, since, until *time.Time) ([]CommitFiled, error) {
	sinceTime, _ := time.Parse("2006", "1999")
	untilTime, _ := time.Parse("2006", "9999")
	if since != nil {
		sinceTime = *since
	}
	if until != nil {
		untilTime = *until
	}

	var query struct {
		Repository struct {
			Ref struct {
				Target struct {
					Commits struct {
						History struct {
							Edges []struct {
								Node struct {
									CommitFiled `graphql:"... on Commit"`
								} `graphql:"node"`
							} `graphql:"edges"`
						} `graphql:"history(first: 100, since: $since, until: $until)"`
					} `graphql:"... on Commit"`
				} `graphql:"target"`
			} `graphql:"ref(qualifiedName: $branch)"`
		} `graphql:"repository(name: $repo,owner: $owner)"`
	}
	params := map[string]interface{}{
		"owner":  githubv4.String(owner),
		"repo":   githubv4.String(repo),
		"branch": githubv4.String(ref),
		"since":  githubv4.GitTimestamp{Time: sinceTime},
		"until":  githubv4.GitTimestamp{Time: untilTime},
	}
	if err := client.client.Query(context.Background(), &query, params); err != nil {
		return nil, err
	}

	result := make([]CommitFiled, 0)
	for _, edge := range query.Repository.Ref.Target.Commits.History.Edges {
		result = append(result, edge.Node.CommitFiled)
	}

	return result, nil
}

// GetCommits history by ref
// History may trace back up to 64 commits due to the limitation of git graphql api
// Ref may be branch or tag or other git refs.
func (client *GithubInfoV4) GetCommitsByRef(owner, repo, ref string) ([]CommitFiled, error) {
	var query struct {
		Repository struct {
			Ref struct {
				Target struct {
					Commits struct {
						History struct {
							Edges []struct {
								Node struct {
									CommitFiled `graphql:"... on Commit"`
								} `graphql:"node"`
							} `graphql:"edges"`
						} `graphql:"history(first: 100)"`
					} `graphql:"... on Commit"`
				} `graphql:"target"`
			} `graphql:"ref(qualifiedName: $branch)"`
		} `graphql:"repository(name: $repo,owner: $owner)"`
	}
	params := map[string]interface{}{
		"owner":  githubv4.String(owner),
		"repo":   githubv4.String(repo),
		"branch": githubv4.String(ref),
	}
	if err := client.client.Query(context.Background(), &query, params); err != nil {
		return nil, err
	}

	result := make([]CommitFiled, 0)
	for _, edge := range query.Repository.Ref.Target.Commits.History.Edges {
		result = append(result, edge.Node.CommitFiled)
	}

	return result, nil
}

// GetCommits history by ref
// History may trace back up to 64 commits due to the limitation of git graphql api
// Ref may be branch or tag or other git refs.
func (client *GithubInfoV4) GetCommitsOfDefaultBranch(owner, repo, startCommitSha string) ([]CommitFiled, PageInfo, error) {
	startCursor := InitClientV4Cursor(startCommitSha)
	var query struct {
		Repository struct {
			Ref struct {
				Target struct {
					Commits struct {
						History struct {
							Edges []struct {
								Node struct {
									CommitFiled `graphql:"... on Commit"`
								} `graphql:"node"`
							} `graphql:"edges"`
							PageInfo `graphql:"pageInfo"`
						} `graphql:"history(first:100,after: $startCommit)"`
					} `graphql:"... on Commit"`
				} `graphql:"target"`
			} `graphql:"defaultBranchRef"`
		} `graphql:"repository(name: $repo,owner: $owner)"`
	}
	params := map[string]interface{}{
		"owner":       githubv4.String(owner),
		"repo":        githubv4.String(repo),
		"startCommit": githubv4.String(startCursor),
	}
	if err := client.client.Query(context.Background(), &query, params); err != nil {
		return nil, PageInfo{}, err
	}

	result := make([]CommitFiled, 0)
	for _, edge := range query.Repository.Ref.Target.Commits.History.Edges {
		result = append(result, edge.Node.CommitFiled)
	}

	return result, query.Repository.Ref.Target.Commits.History.PageInfo, nil
}

// GetCommits history by ref
// History may trace back up to 64 commits due to the limitation of git graphql api
// Ref may be branch or tag or other git refs.
func (client *GithubInfoV4) GetLastCommitOfDefaultBranchUntil(owner, repo string, until time.Time) (*CommitFiled, error) {
	var query struct {
		Repository struct {
			Ref struct {
				Target struct {
					Commits struct {
						History struct {
							Edges []struct {
								Node struct {
									CommitFiled `graphql:"... on Commit"`
								} `graphql:"node"`
							} `graphql:"edges"`
						} `graphql:"history(first:1,until: $until)"`
					} `graphql:"... on Commit"`
				} `graphql:"target"`
			} `graphql:"defaultBranchRef"`
		} `graphql:"repository(name: $repo,owner: $owner)"`
	}

	params := map[string]interface{}{
		"owner": githubv4.String(owner),
		"repo":  githubv4.String(repo),
		"until": githubv4.GitTimestamp{Time: until},
	}
	if err := client.client.Query(context.Background(), &query, params); err != nil {
		return nil, err
	}

	return &query.Repository.Ref.Target.Commits.History.Edges[0].Node.CommitFiled, nil
}

func (client *GithubInfoV4) ClosePullRequestsById(pullRequestID string) error {
	var mutation struct {
		ClosePullRequest struct {
			ClientMutationId githubv4.String `graphql:"clientMutationId"`
		} `graphql:"closePullRequest(input:  $input )"`
	}
	input := githubv4.ClosePullRequestInput{
		PullRequestID: githubv4.ID(pullRequestID),
	}
	params := map[string]interface{}{}

	if err := client.client.Mutate(context.Background(), &mutation, input, params); err != nil {
		return err
	}

	return nil
}

func (client *GithubInfoV4) GetBranchesByCommit(owner, repo, commitId string) ([]string, error) {
	var query struct {
		Repository struct {
			Object struct {
				Commits struct {
					AssociatedPullRequests struct {
						Edges []struct {
							Node struct {
								BaseRef struct {
									Name string `graphql:"name"`
								} `graphql:"baseRef"`
							} `graphql:"node"`
						} `graphql:"edges"`
					} `graphql:"associatedPullRequests(first: 100)"`
				} `graphql:"... on Commit"`
			} `graphql:"object(oid: $commitId)"`
		} `graphql:"repository(name: $repo,owner: $owner)"`
	}

	params := map[string]interface{}{
		"owner":    githubv4.String(owner),
		"repo":     githubv4.String(repo),
		"commitId": githubv4.GitObjectID(commitId),
	}
	if err := client.client.Query(context.Background(), &query, params); err != nil {
		return nil, err
	}

	result := make([]string, 0)
	for _, edges := range query.Repository.Object.Commits.AssociatedPullRequests.Edges {
		result = append(result, edges.Node.BaseRef.Name)
	}

	return result, nil
}

func (client *GithubInfoV4) GetCommitIDByTag(owner, repo, tag string) (string, error) {
	tagName := fmt.Sprintf("refs/tags/%s", tag)
	var query struct {
		Repository struct {
			Ref struct {
				Target struct {
					OID string `graphql:"oid"`
					// Commits struct {
					//
					// } `graphql:"... on Commit"`
				} `graphql:"target"`
			} `graphql:"ref(qualifiedName: $tagName)"`
		} `graphql:"repository(name: $repo,owner: $owner)"`
	}
	params := map[string]interface{}{
		"owner":   githubv4.String(owner),
		"repo":    githubv4.String(repo),
		"tagName": githubv4.String(tagName),
	}

	if err := client.client.Query(context.Background(), &query, params); err != nil {
		return "", err
	}

	return query.Repository.Ref.Target.OID, nil
}
