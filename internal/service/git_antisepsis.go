package service

import (
	"tirelease/commons/git"
	"tirelease/internal/model"
)

func GetUserByGitCode(clientId, clientSecret, code string) (*model.GitUser, error) {
	accessToken, err := git.GetAccessTokenByClient(clientId, clientSecret, code)
	if err != nil {
		return nil, err
	}

	user, err := git.GetUserByToken(accessToken)
	if err != nil {
		return nil, err
	}

	return &model.GitUser{
		GitID:        user.GetID(),
		GitLogin:     user.GetLogin(),
		GitAvatarURL: user.GetAvatarURL(),
		GitName:      user.GetName(),
	}, nil
}
