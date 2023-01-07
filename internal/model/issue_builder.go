package model

import (
	"tirelease/commons/utils"
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

type IssueBuilder struct {
	IssueOption  *entity.IssueOption
	AffectOption *entity.IssueAffectOption

	TriageBuildCommand *TriageBuildCommand
}

type TriageBuildCommand struct {
	NeedTriages bool

	VersionTriages []entity.VersionTriage
}

func (builder IssueBuilder) Option(issueOption *entity.IssueOption, affectOption *entity.IssueAffectOption) IssueBuilder {
	builder.IssueOption = issueOption
	builder.AffectOption = affectOption
	return builder
}

func (builder IssueBuilder) Command(command *TriageBuildCommand) IssueBuilder {
	builder.TriageBuildCommand = command
	return builder
}

func (builder IssueBuilder) BuildArray() ([]Issue, error) {
	option := builder.IssueOption

	command := builder.TriageBuildCommand
	if builder.AffectOption != nil {
		if affects, err := repository.SelectIssueAffect(
			builder.AffectOption,
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

	result, err := builder.buildBareIssues(option)
	if err != nil {
		return nil, err
	}

	if command.NeedTriages {
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

// @deprecated pls use BuildArray function
func (builder IssueBuilder) BuildIssues(option *entity.IssueOption) ([]Issue, error) {
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

func (builder IssueBuilder) buildBareIssues(option *entity.IssueOption) ([]Issue, error) {
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

	result = *fillBlockDefaultValue(&issue, &result)

	return result
}
