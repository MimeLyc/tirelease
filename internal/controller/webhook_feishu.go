package controller

import (
	"net/http"

	"tirelease/internal/controller/feishu_handler"

	"github.com/gin-gonic/gin"
)

func FeishuWebhookHandler(c *gin.Context) {
	handler, err := feishu_handler.GetHandler(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	handler.Handle(c)
}
