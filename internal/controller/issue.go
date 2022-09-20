package controller

import (
	"net/http"
	"tirelease/internal/dto"
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
