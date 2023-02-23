package feishu

import (
	"fmt"
	"testing"
	"tirelease/commons/configs"
	"tirelease/internal/constants"

	"github.com/stretchr/testify/assert"
)

func TestSendTextMsgCard(t *testing.T) {
	configs.LoadConfig("../../config.yaml")
	config := configs.Config

	token, _ := GetAccessTokenFromApp(config.Feishu.AppId, config.Feishu.AppSecret)

	receiver := MsgReceiver{
		IDType: MsgIDTypeEmail,
		ID:     "yuchao.li@pingcap.com",
	}

	msgContent := fmt.Sprintf("You failed to **%s** the **%s**! \n"+
		"The error msg is: \n"+
		"<font color='red'>%s</font>\n"+
		"You can ask the developer of TiRelease for help.\n",
		"watch", "issue",
		`test
            test.watch
        `)

	SendMsgCard(receiver,
		CardMsgWrapper{
			MsgType: "text",
			Msg: NewContentCard("Sorry🙏 !",
				constants.NotifySeverityAlarm,
				[]ContentElement{NewMDCardElement(msgContent)},
			),
		},
		token)
}

func TestSendPostMsgCard(t *testing.T) {

	configs.LoadConfig("../../config.yaml")
	config := configs.Config

	token, _ := GetAccessTokenFromApp(config.Feishu.AppId, config.Feishu.AppSecret)
	receiver := MsgReceiver{
		IDType: MsgIDTypeEmail,
		ID:     "yuchao.li@pingcap.com",
	}

	SendMsgCard(receiver,
		PostMsgWrapper{
			MsgType: "text",
			Msg: ContentPost{
				ZhCnContent: ContentWrapper{
					Title: "Test",
					Rows: [][]ContentElement{
						{
							NewTextContentElement("test11"),
							NewHrefContentElement("https://www.google.com", "google"),
						},
						{

							NewTextContentElement("test21"),
						},
					},
				},
			},
		},
		token)
}

func TestSendButtonCard(t *testing.T) {

	configs.LoadConfig("../../config.yaml")
	config := configs.Config

	token, _ := GetAccessTokenFromApp(config.Feishu.AppId, config.Feishu.AppSecret)
	receiver := MsgReceiver{
		IDType: MsgIDTypeEmail,
		ID:     "yuchao.li@pingcap.com",
	}

	msgContent := fmt.Sprintf("You failed to **%s** the **%s**! \n"+
		"The error msg is: \n"+
		"<font color='red'>%s</font>\n"+
		"You can ask the developer of TiRelease for help.\n",
		"watch", "issue",
		`test
            test.watch
        `)

	err := SendMsgCard(receiver,
		CardMsgWrapper{
			MsgType: "text",
			Msg: NewContentCard("Sorry🙏 !",
				constants.NotifySeverityAlarm,
				[]ContentElement{
					NewMDCardElement(msgContent),
					NewActionElement([]ButtonElement{NewButtonElement("Test Button", nil, InteractiveTypeDanger)}),
				},
			),
		},
		token)
	assert.Nil(t, err)
}
