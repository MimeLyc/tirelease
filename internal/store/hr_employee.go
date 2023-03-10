package store

import (
	"tirelease/internal/entity"
)

func SelectAllHrEmployee() ([]entity.HrEmployee, error) {
	var hrEmployees []entity.HrEmployee
	result := storeGlobalHrEmployeeDB.Find(&hrEmployees)
	if result.Error != nil {
		return nil, result.Error
	}

	return hrEmployees, nil
}
