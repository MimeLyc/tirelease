package service

import (
	"testing"
	"time"

	"tirelease/internal/entity"
	"tirelease/internal/store"
	"tirelease/utils/configs"

	"github.com/stretchr/testify/assert"
)

func TestIssueAffectOperate(t *testing.T) {
	t.Skip()
	// Init
	config := configs.NewConfig(
		"../../"+configs.TestConfig,
		"../../"+configs.TestSecretConfig)
	store.NewStore(config)

	// Init-Data
	var issueAffect = &entity.IssueAffect{
		CreateTime: time.Now(),
		UpdateTime: time.Now(),

		AffectVersion: "5.5.2",
		IssueID:       "100",
		AffectResult:  entity.AffectResultResultUnKnown,
	}
	err := store.CreateOrUpdateIssueAffect(issueAffect)
	assert.Equal(t, true, err == nil)

	// Update
	var updateOption = &entity.IssueAffectUpdateOption{
		IssueID:       "100",
		AffectVersion: "5.5.2",
		AffectResult:  entity.AffectResultResultYes,
	}
	err = IssueAffectOperate(updateOption)
	assert.Equal(t, true, err == nil)
}

func TestIssueAffectOperateWeb(t *testing.T) {
	t.Skip()
	// Init
	config := configs.NewConfig(
		"../../"+configs.TestConfig,
		"../../"+configs.TestSecretConfig)
	store.NewStore(config)

	// Update
	var updateOption = &entity.IssueAffectUpdateOption{
		IssueID:       "I_kwDOAoCpQc5BYBWZ",
		AffectVersion: "5.3",
		AffectResult:  entity.AffectResultResultYes,
	}
	err := IssueAffectOperate(updateOption)
	assert.Equal(t, true, err == nil)
}
