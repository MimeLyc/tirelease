package controller

import (
	"net/http"

	"tirelease/internal/entity"
	"tirelease/internal/store"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func SelectRepo(c *gin.Context) {
	// Params
	option := entity.RepoOption{}
	if err := c.ShouldBindWith(&option, binding.Form); err != nil {
		c.Error(err)
		return
	}

	// Action
	repos, err := store.SelectRepo(&option)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": repos})
}
