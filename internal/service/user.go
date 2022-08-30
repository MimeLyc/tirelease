package service

import "tirelease/internal/model"

func FindUserByCode(clientId, clientSecret, code string) (*model.User, error) {
	user, err := GetUserByGitCode(clientId, clientSecret, code)
	if err != nil {
		return nil, err
	}
	// TODO Replenish user info

	return &model.User{
		GitUser: *user,
	}, nil
}
