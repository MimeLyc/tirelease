package controller

import (
	"net/http"
	"tirelease/internal/dto"
	"tirelease/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func FindSprint(c *gin.Context) {
	request := dto.SprintRequest{}
	if err := c.ShouldBindWith(&request, binding.Form); err != nil {
		c.Error(err)
		return
	}

	// Action
	sprints, err := repository.SelectSprintMetas(&request.SprintOption)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": sprints})

}
