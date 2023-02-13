package model

import (
	"fmt"
	"tirelease/commons/utils"
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

type IssueCmd struct {
	// Filter options
	IssueOption         *entity.IssueOption
	AffectOption        *entity.IssueAffectOption
	VersionTriageOption *entity.VersionTriageOption

	// -- filter by related pr
	ByRelatedPr bool
	PrIDs       []string

	// Field Options
	TriageBuildCommand *TriageBuildCommand
}

type TriageBuildCommand struct {
	WithTriages bool

	VersionTriages []entity.VersionTriage
}

func (cmd IssueCmd) Option(issueOption *entity.IssueOption, affectOption *entity.IssueAffectOption) IssueCmd {
	cmd.IssueOption = issueOption
	cmd.AffectOption = affectOption
	return cmd
}

func (cmd IssueCmd) Command(command *TriageBuildCommand) IssueCmd {
	cmd.TriageBuildCommand = command
	return cmd
}

func (cmd IssueCmd) BuildByNumber(owner, repo string, number int) (*Issue, error) {
	return cmd.buildBareIssue(
		&entity.IssueOption{
			Owner:  owner,
			Repo:   repo,
			Number: number,
		},
	)
}

func (cmd IssueCmd) BuildArray() ([]Issue, error) {
	option := cmd.IssueOption

	command := cmd.TriageBuildCommand
	if cmd.AffectOption != nil {
		if affects, err := repository.SelectIssueAffect(
			cmd.AffectOption,
		); err != nil {
			return nil, err
		} else if len(*affects) == 0 {
			return []Issue{}, nil
		} else {
			affectIssueIds := ExtractIssueIDs(*affects)

			if option.IssueIDs == nil {
				option.IssueIDs = affectIssueIds
			} else {
				option.IssueIDs = utils.Intersects(affectIssueIds, option.IssueIDs)
			}
		}
	}

	result, err := cmd.buildBareIssues(option)
	if err != nil {
		return nil, err
	}

	if command.WithTriages {
		issueIds := extractIssueIdsFromIssueModels(result)
		// TODO: add bellow logic to VersionTriageBuilder
		command, err = command.Triages(
			&entity.VersionTriageOption{
				IssueIDs: issueIds,
			},
		)
		if err != nil {
			return nil, err
		}

		for i := range result {
			issue := &result[i]
			issue.VersionTriages = command.BuildByIssue(issue.Issue)
		}
	}
	return result, nil
}

func (cmd IssueCmd) buildBareIssues(option *entity.IssueOption) ([]Issue, error) {
	issues, err := repository.SelectIssue(
		option,
	)

	if err != nil {
		return nil, err
	}

	ghLogins := extractAuthorGhLoginsFromIssues(issues)
	ghLogins = append(ghLogins, extractAssigneeGhLoginsFromIssues(issues)...)
	userMap, err := UserBuilder{}.BuildUsersByGhLogins(ghLogins)
	if err != nil {
		return nil, err
	}

	result := make([]Issue, 0)
	for _, issue := range *issues {
		issue := issue
		result = append(result, Issue{
			Issue:     issue,
			Assignees: composeAssignees(issue, userMap),
			Author:    userMap[issue.AuthorGHLogin],
		})

	}
	return result, nil
}

func (cmd IssueCmd) buildBareIssue(option *entity.IssueOption) (*Issue, error) {
	issue, err := repository.SelectIssueUnique(
		option,
	)

	if err != nil {
		return nil, err
	}
	if issue == nil {
		return nil, fmt.Errorf("Issue of %v not found.", option)
	}

	ghLogins := extractAuthorGhLoginsFromIssues(&[]entity.Issue{*issue})
	ghLogins = append(ghLogins, extractAssigneeGhLoginsFromIssues(&[]entity.Issue{*issue})...)
	userMap, err := UserBuilder{}.BuildUsersByGhLogins(ghLogins)
	if err != nil {
		return nil, err
	}

	return &Issue{
		Issue:     *issue,
		Assignees: composeAssignees(*issue, userMap),
		Author:    userMap[issue.AuthorGHLogin],
	}, nil
}

func (command *TriageBuildCommand) Triages(option *entity.VersionTriageOption) (*TriageBuildCommand, error) {
	triages, err := repository.SelectVersionTriage(option)
	if err != nil {
		return command, err
	}

	command.VersionTriages = *triages

	return command, nil
}

func (command *TriageBuildCommand) BuildByIssue(issue entity.Issue) []entity.VersionTriage {
	result := make([]entity.VersionTriage, 0)
	for _, triage := range command.VersionTriages {
		triage := triage
		if triage.IssueID == issue.IssueID {
			result = append(result, triage)
		}
	}

	return result
}
