package git

import (
	"context"
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

// GetCommits history by ref
// History may trace back up to 64 commits due to the limitation of git graphql api
// Ref may be branch or tag or other git refs.
func (client *GithubInfoV4) GetCommitsByRef(owner, repo, ref string, since, until *time.Time) ([]CommitFiled, error) {
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
						} `graphql:"history(since: $since, until: $until)"`
					} `graphql:"... on Commit"`
				} `graphqjkjl:"target"`
				// Node []struct {
				// 	Name   string `graphql:"name"`
				// 	Target struct {
				// 		commits CommitFiled `graphql:"... on Commit"`
				// 	} `graphql:"target"`
				// } `graphql:"nodes"`
				// } `graphql:"ref(first: 50, refPrefix: \"refs/heads/\", query: $branch, orderBy: {field: TAG_COMMIT_DATE, direction: DESC})"`
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
