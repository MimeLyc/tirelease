package feishu_handler

import (
	"fmt"

	"github.com/gin-gonic/gin"

	. "tirelease/commons/log"
)

type EventType string

const (
	TypeReceiveMsgV1 EventType = "im.message.receive_v1"
	TypeVerifyCode   EventType = "url_verification"
)

type Handler interface {
	shouldHandle(c *gin.Context) (bool, error)
	Handle(c *gin.Context)
}

var Handlers = make([]Handler, 0)

func init() {
	Handlers = append(Handlers, &UrlVerifyHandler{})
	Handlers = append(Handlers, &MsgReceiverHandler{})
	Handlers = append(Handlers, &ActionReceiverHandler{})
}

func GetHandler(c *gin.Context) (Handler, error) {
	for _, handler := range Handlers {
		shouldHandler, err := handler.shouldHandle(c)
		if err != nil {
			Log.Errorf(err, "ShouldHandler valified error")
			continue
		}

		if shouldHandler {
			return handler, nil
		}
	}

	return nil, fmt.Errorf("Couldn't find handler for %v", c)
}
