package notify

import "tirelease/commons/feishu"

func sendFeishuMsgByEmail(email string, msg feishu.MsgWrapper) error {
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

	contentCard := feishu.NewContentCard(
		content.Header,
		content.Severity,
		content.ParseToFeishuContent(),
	)

	cardWrapper := feishu.CardMsgWrapper{
		MsgType: "interactive",
		Msg:     contentCard,
	}
	return sendFeishuMsgByEmail(email, cardWrapper)
}

func SendFeishuFormattedByGroup(group string, content NotifyContent) error {
	token, err := feishu.GetAccessToken()
	if err != nil {
		return err
	}

	contentCard := feishu.NewContentCard(
		content.Header,
		content.Severity,
		content.ParseToFeishuContent(),
	)

	cardWrapper := feishu.CardMsgWrapper{
		MsgType: "interactive",
		Msg:     contentCard,
	}

	receiver := feishu.MsgReceiver{
		IDType: feishu.MsgIDTypeChat,
		ID:     group,
	}

	return feishu.SendMsgCard(receiver, cardWrapper, token)
}

func ReplyFeishuByMessageId(messageId string, content NotifyContent) error {
	token, err := feishu.GetAccessToken()
	if err != nil {
		return err
	}

	contentCard := feishu.NewContentCard(
		content.Header,
		content.Severity,
		content.ParseToFeishuContent(),
	)

	cardWrapper := feishu.CardMsgWrapper{
		MsgType: "interactive",
		Msg:     contentCard,
	}

	return feishu.ReplyMessage(messageId, cardWrapper, token)
}
