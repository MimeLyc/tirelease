package feishu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const replyApi = "https://open.feishu.cn/open-apis/im/v1/messages/%s/reply"

type ReplyReq struct {
	MsgType string `json:"msg_type"`
	Content string `json:"content"`
}

func ReplyMessage(messageId string, msg MsgWrapper, token string) error {
	requestUrl := fmt.Sprintf(replyApi, messageId)
	reqBody := ReplyReq{
		MsgType: msg.GetMsgType(),
		Content: msg.GetMsgJson(),
	}

	jsonContent, err := json.Marshal(reqBody)

	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", requestUrl, bytes.NewBuffer(jsonContent))
	if err != nil {
		return err
	}
	authString := fmt.Sprintf("Bearer %s", token)
	req.Header.Set("Authorization", authString)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}
	result := MsgResp{}
	json.Unmarshal(body, &result)

	if result.Code != 0 {
		return fmt.Errorf("Reply Feishu message %s error: error code %d, error msg: %s", messageId, result.Code, result.Message)
	}

	return nil

}
