package feishu_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type UrlVerifier struct {
	Challenge string `json:"challenge,omitempty" form:"challenge"`
	Token     string `json:"token" form:"token"`
	Type      string `json:"type" form:"type"`
}

type UrlVerifyHandler struct {
}

func (h *UrlVerifyHandler) shouldHandle(c *gin.Context) (bool, error) {
	verifier := UrlVerifier{}

	if err := c.ShouldBindBodyWith(&verifier, binding.JSON); err != nil {
		return false, err
	}

	return verifier.Type == string(TypeVerifyCode), nil
}

func (h *UrlVerifyHandler) Handle(c *gin.Context) {
	verifier := UrlVerifier{}

	if err := c.ShouldBindBodyWith(&verifier, binding.JSON); err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, gin.H{"challenge": verifier.Challenge})
}
