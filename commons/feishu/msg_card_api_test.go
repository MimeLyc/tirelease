package feishu

import (
	"testing"
	"tirelease/commons/configs"
)

func TestSendTextMsgCard(t *testing.T) {
	configs.LoadConfig("../../config.yaml")
	config := configs.Config

	token, _ := GetAccessTokenFromApp(config.Feishu.AppId, config.Feishu.AppSecret)

	receiver := MsgReceiver{
		IDType: MsgIDTypeEmail,
		ID:     "yuchao.li@pingcap.com",
	}

	msgContent := "{\"text\": \"this is a Feishu Msg test from TiRelease. \\n this Test new line. test test \"}"

	SendMsgCard(receiver,
		TextMsgWrapper{
			MsgType: "text",
			Msg:     msgContent,
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
