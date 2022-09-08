package model

import "tirelease/internal/entity"

type SprintMeta struct {
	entity.SprintMeta
}

func (sprint SprintMeta) GetMajorVersion() int {
	major, _, _, _ := ComposeVersionAtom(sprint.MinorVersionName)
	return major
}

func (sprint SprintMeta) GetMinorVersion() int {
	_, minor, _, _ := ComposeVersionAtom(sprint.MinorVersionName)
	return minor
}
