package service

import (
	"tirelease/internal/entity"
	"tirelease/internal/model"
)

func FindUserByCode(clientId, clientSecret, code string) (*model.User, error) {
	user, err := model.GetUserByGitCode(clientId, clientSecret, code)
	if err != nil {
		return nil, err
	}
	// TODO Replenish user info

	return user, nil
}

func FindEmployees(options entity.EmployeeOptions) ([]model.User, error) {
	employees, err := model.UserCmd{
		Options: &options,
	}.BuildEmployees()

	if err != nil {
		return nil, err
	}

	return employees, nil
}
