package store

import (
	"fmt"
	"time"

	"tirelease/internal/entity"

	"github.com/pkg/errors"
)

// Implement
func TriageItemInsert(triageItem *entity.TriageItem) error {
	triageItem.CreateTime = time.Now()
	triageItem.UpdateTime = time.Now()

	if err := tempDB.DB.Create(&triageItem).Error; err != nil {
		return errors.Wrap(err, fmt.Sprintf("create triage item: %+v failed", triageItem))
	}
	return nil
}

func TriageItemSelect(option *entity.TriageItemOption) (*[]entity.TriageItem, error) {
	var triageItems []entity.TriageItem

	if err := tempDB.DB.Where(option).Find(&triageItems).Error; err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("find triage item: %+v failed", option))
	}
	return &triageItems, nil
}
