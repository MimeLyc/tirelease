package model

import (
	"tirelease/internal/entity"
)

func ParseFromEntityPickTriage(status entity.VersionTriageResult) StateText {
	return StateText(status)
}

func ParseToEntityPickTriage(status StateText) entity.VersionTriageResult {
	return entity.VersionTriageResult(status)
}
