package controller

import (
	"net/http"

	"tirelease/internal/dto"
	"tirelease/internal/entity"
	"tirelease/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func CreateOrUpdateVersionTriage(c *gin.Context) {
	// Params
	versionTriageModifyOption := entity.VersionTriageModifyOption{}
	if err := c.ShouldBindWith(&versionTriageModifyOption, binding.JSON); err != nil {
		c.Error(err)
		return
	}

	// Action
	versionTriageInfo, err := service.SaveVersionTriageInfo(&versionTriageModifyOption.VersionTriage, versionTriageModifyOption.UpdatedVars...)
	if nil != err {
		c.Error(err)
		return
	}

	// Response
	var statusCode int = 200
	if nil != versionTriageInfo && versionTriageInfo.IsFrozen && versionTriageInfo.IsAccept {
		statusCode = 202
	}
	c.JSON(statusCode, gin.H{"status": "ok", "data": versionTriageInfo})
}

func CreateOrUpdateBatchTriages(c *gin.Context) {
	// Params
	option := make([]entity.VersionTriageOption, 0)

	if err := c.ShouldBindWith(&option, binding.JSON); err != nil {
		c.Error(err)
		return
	}

	// Action
	err := service.CreateOrUpdateIssueTriages(&option)
	if nil != err {
		c.Error(err)
		return
	}

	//		Response
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func SelectVersionTriageInfo(c *gin.Context) {
	// Params
	versionTriageInfoQuery := dto.VersionTriageInfoQuery{}
	if err := c.ShouldBindUri(&versionTriageInfoQuery); err != nil {
		c.Error(err)
		return
	}
	if versionTriageInfoQuery.Version != "" {
		versionTriageInfoQuery.VersionName = versionTriageInfoQuery.Version
	}

	// Action
	versionTriageInfos, response, err := service.FindVersionTriageInfo(&versionTriageInfoQuery)
	if nil != err {
		c.Error(err)
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{"data": versionTriageInfos, "response": response})
}

func SelectVersionTriageResult(c *gin.Context) {
	var enumResult = struct {
		VersionTriageResultUnKnown      entity.VersionTriageResult
		VersionTriageResultAccept       entity.VersionTriageResult
		VersionTriageResultAcceptFrozen entity.VersionTriageResult
		VersionTriageResultLater        entity.VersionTriageResult
		VersionTriageResultWontFix      entity.VersionTriageResult
		VersionTriageResultReleased     entity.VersionTriageResult
	}{
		VersionTriageResultUnKnown:      entity.VersionTriageResultUnKnown,
		VersionTriageResultAccept:       entity.VersionTriageResultAccept,
		VersionTriageResultAcceptFrozen: entity.VersionTriageResultAcceptFrozen,
		VersionTriageResultLater:        entity.VersionTriageResultLater,
		VersionTriageResultWontFix:      entity.VersionTriageResultWontFix,
		VersionTriageResultReleased:     entity.VersionTriageResultReleased,
	}

	c.JSON(http.StatusOK, gin.H{"data": enumResult})
}
