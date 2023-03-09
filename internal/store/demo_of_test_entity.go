package store

import (
	"fmt"
	"time"

	"tirelease/internal/entity"

	"github.com/pkg/errors"
)

// Implement
func TestEntityInsert(testEntity *entity.TestEntity) error {
	if testEntity.CreateTime.IsZero() {
		testEntity.CreateTime = time.Now()
	}
	if testEntity.UpdateTime.IsZero() {
		testEntity.UpdateTime = time.Now()
	}

	if err := tempDB.DB.Create(&testEntity).Error; err != nil {
		return errors.Wrap(err, fmt.Sprintf("create test entity: %+v failed", testEntity))
	}
	return nil
}

func TestEntitySelect(option *entity.TestEntityOption) (*[]entity.TestEntity, error) {
	var testEntities []entity.TestEntity

	if err := tempDB.DB.Where(option).Find(&testEntities).Error; err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("find test entity: %+v failed", option))
	}
	return &testEntities, nil
}
