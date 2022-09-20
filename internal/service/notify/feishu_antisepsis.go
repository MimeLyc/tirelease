package notify

import "tirelease/commons/feishu"

func sendFeishuPostMsgByEmail(email string, msg feishu.PostMsgWrapper) error {
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

func SendFeishuFormattedByEmail(email string, content NotifyContent) error {
	contentPost := feishu.ContentPost{
		ZhCnContent: feishu.ContentWrapper{
			Title: content.Header,
			Rows:  content.ParseToFeishuContent(),
		},
	}

	postWrapper := feishu.PostMsgWrapper{
		MsgType: "text",
		Msg:     contentPost,
	}
	return sendFeishuPostMsgByEmail(email, postWrapper)
}
