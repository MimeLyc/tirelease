package feishu_handler

import (
	"fmt"
	"net/http"

	. "tirelease/commons/log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type ActionReceiverHandler struct {
}

func (h *ActionReceiverHandler) shouldHandle(c *gin.Context) (bool, error) {
	verifier := ActionReceive{}

	if err := c.ShouldBindBodyWith(&verifier, binding.JSON); err != nil {
		return false, err
	}

	return verifier.Action != nil &&
			verifier.Action.Value != nil,
		nil
}

func (h *ActionReceiverHandler) Handle(c *gin.Context) {
	receiver := ActionReceive{}

	if err := c.ShouldBindBodyWith(&receiver, binding.JSON); err != nil {
		Log.Errorf(err, "Bind request with ActionReceiver error")
		c.JSON(http.StatusInternalServerError, err)
	}

	Log.Infof("Handle feishu action: %v", receiver)
	if err := deliverAction(receiver); err != nil {
		Log.Errorf(err, "DeliverAction err, receive: %v", receiver)
		// Return OK to feishu, so that feishu will not retry
		// The actual error is logged in response of msg
		c.JSON(http.StatusOK, nil)
	}
	c.JSON(http.StatusOK, nil)
}

type ActionReceive struct {
	OpenID    string       `json:"open_id"`
	UserID    string       `json:"user_id"`
	OpenMsgID string       `json:"open_message_id"`
	Action    *ActionValue `json:"action"`
}

type ActionValue struct {
	Value map[string]interface{} `json:"value"`
	Tag   string                 `json:"tag"`
}

func (receive *ActionReceive) Validate() error {
	_, ok1 := receive.Action.Value["register_object"]
	_, ok2 := receive.Action.Value["register_object_id"]
	_, ok3 := receive.Action.Value["register_action"]
	if !(ok1 && ok2 && ok3) {
		return fmt.Errorf("Invalid ActionReceive: %v", receive)
	}

	return nil
}
