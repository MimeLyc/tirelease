package feishu_handler

import (
	"net/http"
	"strings"

	. "tirelease/commons/log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type MsgReceiverHandler struct {
}

func (h *MsgReceiverHandler) shouldHandle(c *gin.Context) (bool, error) {
	verifier := MsgReceiveV1{}

	if err := c.ShouldBindBodyWith(&verifier, binding.JSON); err != nil {
		return false, err
	}

	return verifier.Schema.verifyVersion() &&
			verifier.Header.EventType == string(TypeReceiveMsgV1) &&
			!strings.Contains(verifier.Event.Message.Content, "@_all"),
		nil
}

func (msgreceiverhandler *MsgReceiverHandler) Handle(c *gin.Context) {
	receiver := MsgReceiveV1{}

	if err := c.ShouldBindBodyWith(&receiver, binding.JSON); err != nil {
		Log.Errorf(err, "Bind request with MsgReceiver error")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	Log.Infof("Handle feishu message: %v", receiver)
	if err := deliverCmd(receiver); err != nil {
		Log.Errorf(err, "DeliverCmd err, receive: %v", receiver)
		// Return OK to feishu, so that feishu will not retry
		// The actual error is logged in response of msg
		c.JSON(http.StatusOK, nil)
		return
	}
	c.JSON(http.StatusOK, nil)
}

type MsgReceiveV1 struct {
	Schema
	Header msgReceiveV1Header `json:"header,omitempty" form:"header"`
	Event  msgReceiveV1Event  `json:"event,omitempty"`
}

type msgReceiveV1Header struct {
	EventId   string `json:"event_id,omitempty"`
	EventType string `json:"event_type,omitempty"`
	Token     string `json:"token,omitempty"`
}

type msgReceiveV1Event struct {
	Sender  msgReceiveV1Sender `json:"sender,omitempty"`
	Message msgReceiveV1Msg    `json:"message,omitempty"`
}

type msgReceiveV1Sender struct {
	SenderID msgReceiveV1SenderId `json:"sender_id,omitempty"`
}

type msgReceiveV1SenderId struct {
	OpenID string `json:"open_id,omitempty"`
}

type msgReceiveV1Msg struct {
	MessageID string                `json:"message_id,omitempty"`
	ChatID    string                `json:"chat_id,omitempty"`
	ChatType  string                `json:"chat_type,omitempty"`
	Content   string                `json:"content,omitempty"`
	Mentions  []msgReceiveV1Mention `json:"mentions,omitempty"`
}

type msgReceiveV1Mention struct {
	Key  string `json:"key,omitempty"`
	Name string `json:"name,omitempty"`
}
