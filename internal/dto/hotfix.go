package dto

import "tirelease/internal/model"

type HotfixSaveRequest struct {
	model.Hotfix
	OperatorEmail string `json:"operator_email"`
}
