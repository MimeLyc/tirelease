package model

import "tirelease/internal/entity"

func ParseFromEntityBlockTriage(status entity.BlockVersionReleaseResult) StateText {
	return StateText(status)
}

func ParseToEntityBlockTriage(status StateText) entity.BlockVersionReleaseResult {
	return entity.BlockVersionReleaseResult(status)
}
