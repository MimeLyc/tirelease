package store

import (
	"fmt"
	"tirelease/internal/entity"

	"github.com/pkg/errors"
	"gorm.io/gorm/clause"
)

func BatchCreateOrUpdateEmployees(employees []entity.Employee) error {
	if err := storeGlobalDB.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&employees).Error; err != nil {
		return errors.Wrap(err, fmt.Sprintf("create or update employees: %+v failed", employees))
	}
	return nil
}

func BatchSelectEmployeesByGhLogins(githubLogins []string) ([]entity.Employee, error) {
	employees := make([]entity.Employee, 0)
	if err := storeGlobalDB.DB.Where("github_id in ?", githubLogins).Find(&employees).Error; err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("select users by github logins: %+v failed", githubLogins))
	}
	return employees, nil
}

func BatchSelectEmployeesByEmails(emails []string) ([]entity.Employee, error) {
	employees := make([]entity.Employee, 0)
	if err := storeGlobalDB.DB.Where("email in ?", emails).Find(&employees).Error; err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("select users by emails: %+v failed", emails))
	}
	return employees, nil
}
