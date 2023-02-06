package model

import (
	"tirelease/commons/git"
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

// SelectIssuePrRelationsByVersion function select and compose the issues with their related PRs.
// @major, @minor: the major and minor part of target version
// @option: the issue option to filter issues
//
//	@limitAffect: there are cases that issues not (labeled) affect the version but have cherry-pick on the version branch. \
//		And such cases need to be compatible while analyzing triage infos.
//
//	return \
//	  empty array while there is no issue. \
//			Compose issue with empty array while there is no related  PRs.
func SelectIssuePrRelationsByVersion(major, minor int, option entity.IssueOption, limitAffect bool) ([]IssuePrRelation, error) {
	versionName := ComposeVersionMinorNameByNumber(major, minor)
	branchName := git.ReleaseBranchPrefix + versionName
	if limitAffect {
		if affects, err := repository.SelectIssueAffect(
			&entity.IssueAffectOption{
				AffectVersion: versionName,
				AffectResult:  entity.AffectResultResultYes,
				IssueIDs:      option.IssueIDs,
			},
		); err != nil {
			return nil, err
		} else if len(*affects) == 0 {
			return []IssuePrRelation{}, nil
		} else {
			issueIds := ExtractIssueIDs(*affects)
			option.IssueIDs = issueIds
		}
	}

	issues, err := repository.SelectIssue(&option)
	if err != nil {
		return nil, err
	}

	issueIDs := extractIssueIdsFromIssues(*issues)
	issuePrRelations, err := repository.SelectIssuePrRelation(
		&entity.IssuePrRelationOption{
			IssueIDs: issueIDs,
		},
	)
	if err != nil {
		return nil, err
	}

	prids := extractPrIdsFromIssuePrRelation(*issuePrRelations)
	prs, err := PullRequestCmd{
		PROptions: &entity.PullRequestOption{
			BaseBranch:     branchName,
			PullRequestIDs: prids,
		},
	}.Build()

	result := make([]IssuePrRelation, 0)

	for _, issue := range *issues {
		issue := issue
		issuePrs := getPRsByIssueRelation(*issuePrRelations, issue.IssueID, &prs)
		if err != nil {
			return nil, err
		}

		result = append(result, IssuePrRelation{
			Major:      major,
			Minor:      minor,
			Issue:      &issue,
			RelatedPrs: issuePrs,
		})
	}

	return result, nil
}

func SelectIssuePrRelationByIds(issueIDs []string) ([]entity.IssuePrRelation, error) {
	issuePrRelationAll := make([]entity.IssuePrRelation, 0)

	if len(issueIDs) > 0 {
		issuePrRelationAlls, err := repository.SelectIssuePrRelation(
			&entity.IssuePrRelationOption{
				IssueIDs: issueIDs,
			},
		)

		if nil != err {
			return nil, err
		}
		issuePrRelationAll = append(issuePrRelationAll, (*issuePrRelationAlls)...)
	}

	return issuePrRelationAll, nil
}
