package service

import (
	"tirelease/internal/entity"
)

func extractIssueIDsFromTriage(triages []entity.VersionTriage) []string {
	issueIDs := make([]string, 0)
	for _, triage := range triages {
		issueIDs = append(issueIDs, triage.IssueID)
	}
	return issueIDs
}

func filterPRsByBranch(prs []entity.PullRequest, branch string) []entity.PullRequest {
	var filteredPRs []entity.PullRequest
	for _, pr := range prs {
		if pr.BaseBranch == branch {
			filteredPRs = append(filteredPRs, pr)
		}
	}
	return filteredPRs
}

// Notes: the version of pullrequests in the params is the same with the triage
// TODO: refactor the model of version triage to contain the related info.
func mapVersionTriagesWithPrs(triages []entity.VersionTriage, issuePrRelation []entity.IssuePrRelation, prs []entity.PullRequest) map[entity.VersionTriage][]entity.PullRequest {
	triagePRMap := make(map[entity.VersionTriage][]entity.PullRequest)
	for _, triage := range triages {
		for _, relation := range issuePrRelation {
			if triage.IssueID != relation.IssueID {
				continue
			}

			for _, pr := range prs {
				if relation.PullRequestID == pr.PullRequestID {
					triagePRMap[triage] = append(triagePRMap[triage], pr)
				}
			}
		}
	}
	return triagePRMap
}
