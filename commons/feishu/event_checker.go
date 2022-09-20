package feishu

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type UrlVerifier struct {
	Challenge string `json:"challenge,omitempty" form:"challenge"`
	Token     string `json:"token" form:"token"`
	Type      string `json:"type" form:"type"`
}

const verifyCode string = "url_verification"

func VerifyWebhookSetRequest(c *gin.Context) (UrlVerifier, bool) {
	verifier := UrlVerifier{}

	if err := c.ShouldBindWith(&verifier, binding.JSON); err != nil {
		c.Error(err)
	}
	return verifier, verifier.Type == verifyCode
}
