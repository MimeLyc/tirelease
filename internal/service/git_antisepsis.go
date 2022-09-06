package service

import (
	"time"
	"tirelease/commons/git"
	"tirelease/internal/model"
)

type GitCommit struct {
	Oid            string
	AbbreviatedOid string
	CommittedTime  time.Time
	PushedTime     time.Time
}

type GitBranch struct {
	FirstCommit   GitCommit
	Owner         string
	Repo          string
	Name          string
	QualifiedName string
	PushedTime    time.Time
}

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
