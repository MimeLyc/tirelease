package constants

type NotifyTriggerType string

const (
	NotifyTriggerVersion NotifyTriggerType = "Release"
	NotifyTriggerAction  NotifyTriggerType = "Action"
	NotifyTriggerCron    NotifyTriggerType = "Cron"
)

type EventRegisterPlatform string

const EventRegisterPlatformFeishu EventRegisterPlatform = "feishu"

// Mainly related to feishu chat type
// see: [接收消息 \- 服务端文档 \- 开发文档 \- 飞书开放平台](https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/im-v1/message/events/receive)
type EventRegisterUserType string

const (
	EventRegisterGroup EventRegisterUserType = "group"
	EventRegisterP2P   EventRegisterUserType = "p2p"
	EventRegisterTopic EventRegisterUserType = "topic"
)

type EventRegisterAction string

const (
	EventRegisterFatFrozenApprove EventRegisterAction = "frozen-approve-fat"
	EventRegisterApprove          EventRegisterAction = "approve"
	EventRegisterDeny             EventRegisterAction = "deny"
	EventRegisterPendingApproval  EventRegisterAction = "pending-approval"
)

type EventRegisterObject string

const (
	EventRegisterObjectIssue  EventRegisterObject = "issue"
	EventRegisterObjectHotfix EventRegisterObject = "hotfix"
)

type NotifySeverity string

const (
	NotifySeverityInfo  NotifySeverity = "INFO"
	NotifySeverityAlarm NotifySeverity = "ALARM"
)
