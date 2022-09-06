package service

import "tirelease/commons/feishu"

func SendFeishuPostMsgByEmail(email string, msg feishu.PostMsgWrapper) error {
	token, err := feishu.GetAccessToken()
	if err != nil {
		return err
	}

	receiver := feishu.MsgReceiver{
		IDType: feishu.MsgIDTypeEmail,
		ID:     email,
	}
	return feishu.SendMsgCard(receiver, msg, token)

}
