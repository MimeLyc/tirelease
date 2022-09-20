package dto

type SprintIssueNotificationRequest struct {
	*SprintMetaRequest
	Email string `json:"email" form:"email"`
}
