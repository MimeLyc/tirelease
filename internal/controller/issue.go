package controller

import (
	"net/http"
	"tirelease/internal/dto"
	"tirelease/internal/entity"
	"tirelease/internal/model"
	"tirelease/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func NotifySprintIssueInfo(c *gin.Context) {
	// Params
	option := dto.SprintIssueNotificationRequest{}

	if err := c.ShouldBindWith(&option, binding.Form); err != nil {
		c.Error(err)
		return
	}

	// Action
	err := service.NotifySprintBugMetrics(*option.Major, *option.Minor, option.Email)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func FindSingleIssueInfo(c *gin.Context) {
	// Params
	option := dto.IssueRelationInfoQuery{}

	if err := c.ShouldBindUri(&option); err != nil {
		c.Error(err)
		return
	}
	if option.Page == 0 {
		option.Page = 1
	}
	if option.PerPage == 0 {
		option.PerPage = 10
	}
	option.ParamFill()

	// Action
	issue, err := model.SelectIssueTriage(option.IssueID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": mapToIssueTriageDTO(*issue)})
}

func mapToIssueTriageDTO(issue model.IssueTriage) dto.IssueTriage {
	return dto.IssueTriage{
		Issue:          &issue.Issue,
		MasterPrs:      &issue.MasterPrs,
		VersionTriages: mapToVersionTriages(*issue.VersionTriages),
	}
}

func mapToVersionTriages(triages []model.IssueVersionTriage) *[]dto.VersionTriage {
	result := make([]dto.VersionTriage, 0)
	for _, triage := range triages {
		triage := triage
		result = append(result, dto.VersionTriage{
			ReleaseVersion:    triage.Version.ReleaseVersion,
			VersionPrs:        &triage.RelatedPrs,
			PickTriageResult:  model.ParseToEntityPickTriage(triage.PickTriage.State.StateText),
			BlockTriageResult: model.ParseToEntityBlockTriage(triage.BlockTriage.State.StateText),
			IsBlock:           model.ParseToEntityPickTriage(triage.BlockTriage.State.StateText) == entity.VersionTriageResult(entity.BlockVersionReleaseResultBlock),
			AffectResult:      triage.Affect,
			IsAffect:          triage.Affect == entity.AffectResultResultYes,

			Comment:     triage.Entity.Comment,
			ChangedItem: triage.Entity.ChangedItem,
			MergeStatus: triage.GetMergeStatus(),
		})
	}

	return &result
}
