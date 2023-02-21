package model

import (
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

type User struct {
	// Basic Info
	Name         string `json:"name"`
	Email        string `json:"email"`
	IsActive     bool   `json:"active"`
	JobNumber    string `json:"job_number"`
	IsEmployee   bool   `json:"is_employ"`
	HrEmployeeID string `json:"hr_id"`

	// Git Info
	GitUser
}

type GitUser struct {
	GitID        int64  `json:"git_id"`
	GitLogin     string `json:"git_login"`
	GitAvatarURL string `json:"git_avatar_url"`
	GitName      string `json:"git_name"`
}

type UserBuilder struct {
}

func (builder UserBuilder) BuildByGhLogin(login string) (*User, error) {
	employees, err := repository.BatchSelectEmployeesByGhLogins([]string{login})
	if err != nil {
		return nil, err
	}

	if len(employees) == 0 {
		return &User{
			IsEmployee: false,
			GitUser: GitUser{
				GitLogin: login,
			},
		}, nil
	}

	employee := employees[0]
	user := mapEmployeeEntity2User(employee)

	return &user, nil
}

// Select Employee info from db
// Compose TiRelease User info,
// only return Github infos while the user is not employee of PingCAP
func (builder UserBuilder) BuildUsersByGhLogins(logins []string) (map[string]User, error) {
	result := make(map[string]User)
	employees, err := repository.BatchSelectEmployeesByGhLogins(logins)
	if err != nil {
		return nil, err
	}

	for _, employee := range employees {
		employee := employee
		user := mapEmployeeEntity2User(employee)
		result[user.GitLogin] = user
	}

	for _, ghId := range logins {
		if _, ok := result[ghId]; !ok {
			result[ghId] = User{
				IsEmployee: false,
				GitUser: GitUser{
					GitLogin: ghId,
				},
			}
		}
	}

	return result, nil
}

func mapEmployeeEntity2User(employee entity.Employee) User {
	return User{
		Name:         employee.Name,
		Email:        employee.Email,
		IsActive:     employee.IsActive,
		IsEmployee:   true,
		HrEmployeeID: employee.HrEmployeeID,
		GitUser: GitUser{
			GitLogin: employee.GithubId,
			GitName:  employee.GhName,
		},
	}
}
