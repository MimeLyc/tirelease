package controller

import (
	"net/http"
	"tirelease/internal/entity"
	"tirelease/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func FindUser(c *gin.Context) {
	// Params
	option := UserRequest{}
	if err := c.ShouldBindQuery(&option); err != nil {
		c.Error(err)
		return
	}

	// Action
	user, err := HandleFindUser(option)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// Find employees of PingCAP
func FindEmployees(c *gin.Context) {
	// Params
	option := entity.EmployeeOptions{}
	if err := c.ShouldBindWith(&option, binding.Form); err != nil {
		c.Error(err)
		return
	}

	// Action
	employees, err := service.FindEmployees(option)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": employees})
}
