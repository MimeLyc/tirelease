package controller

import (
	"net/http"

	"tirelease/commons/feishu"

	"github.com/gin-gonic/gin"
)

func FeishuWebhookHandler(c *gin.Context) {
	// verify feishu event setup webhook.
	urlVerifier, ok := feishu.VerifyWebhookSetRequest(c)
	if ok {
		c.JSON(http.StatusOK, gin.H{"challenge": urlVerifier.Challenge})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
