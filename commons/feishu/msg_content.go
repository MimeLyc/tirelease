package feishu

import "encoding/json"

type MsgWrapper interface {
	GetMsgType() string
	GetMsgJson() string
}

type TextMsgWrapper struct {
	MsgType string `json:"msg_type"`
	Msg     string `json:"content"`
}

func (msg TextMsgWrapper) GetMsgType() string {
	return "text"
}

func (msg TextMsgWrapper) GetMsgJson() string {
	return msg.Msg
}

var _ MsgWrapper = TextMsgWrapper{}

type PostMsgWrapper struct {
	MsgType string      `json:"msg_type"`
	Msg     ContentPost `json:"content"`
}

func (msg PostMsgWrapper) GetMsgType() string {
	return "post"
}

func (msg PostMsgWrapper) GetMsgJson() string {
	jsonData, _ := json.Marshal(msg.Msg)
	return string(jsonData)
}

var _ MsgWrapper = PostMsgWrapper{}

type ContentPost struct {
	ZhCnContent ContentWrapper `json:"zh_cn"`
	// Post PostCN `json:"post,omitempty"`
}

type PostCN struct {
	ZhCnContent ContentWrapper `json:"zh_cn"`
}

type ContentWrapper struct {
	Title string             `json:"title"`
	Rows  [][]ContentElement `json:"content"`
}

type ContentElement interface {
}

type TextContentElement struct {
	Tag      string `json:"tag" default:"text"`
	UnEscape bool   `json:"un_escape,omitempty" default:"false"`
	Text     string `json:"text"`
}

func NewTextContentElement(text string) TextContentElement {
	return TextContentElement{
		Tag:  "text",
		Text: text,
	}
}

var _ ContentElement = TextContentElement{}

type HrefContentElement struct {
	Tag      string `json:"tag" default:"a"`
	UnEscape bool   `json:"un_escape,omitempty" default:"false"`
	Href     string `json:"href"`
	Text     string `json:"text"`
}

func NewHrefContentElement(href, text string) HrefContentElement {
	return HrefContentElement{
		Tag:  "a",
		Href: href,
		Text: text,
	}
}

var _ ContentElement = HrefContentElement{}
