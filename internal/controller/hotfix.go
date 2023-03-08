package controller

import (
	"net/http"
	"tirelease/internal/dto"
	"tirelease/internal/entity"
	"tirelease/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func HandleSaveHotfix(c *gin.Context) {
	// Params
	hotfix := dto.HotfixSaveRequest{}
	if err := c.ShouldBindWith(&hotfix, binding.JSON); err != nil {
		c.Error(err)
		return
	}

	// Action
	respHotfix, err := service.SaveHotfix(hotfix)
	if nil != err {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK,
		gin.H{
			"data":   respHotfix,
			"status": "ok",
		},
	)
}

func HandleFindHotfix(c *gin.Context) {
	// Params
	option := entity.HotfixOptions{}
	if err := c.ShouldBindWith(&option, binding.Form); err != nil {
		c.Error(err)
		return
	}

	// Action
	hotfixes, err := service.FindHotfixes(option)
	if nil != err {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": hotfixes})
}

func HandleFindSingleHotfix(c *gin.Context) {
	// Params
	option := entity.HotfixOptions{}
	if err := c.ShouldBindUri(&option); err != nil {
		c.Error(err)
		return
	}

	if len(option.Name) == 0 {
		return
	}

	// Action
	hotfix, err := service.FindHotfixByName(option.Name)
	if nil != err {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": hotfix})
}
