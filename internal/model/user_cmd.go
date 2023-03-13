package model

import (
	"tirelease/internal/entity"
	"tirelease/internal/store"
)

type UserCmd struct {
	Options *entity.EmployeeOptions

	// TODO: Refactor with this tags to make the users who is not employee
	// MustBeEmployee bool
}

func (builder UserCmd) BuildByGhLogin(login string) (*User, error) {
	employees, err := store.BatchSelectEmployeesByGhLogins([]string{login})
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

func (builder UserCmd) BuildByEmail(email string) (*User, error) {
	employees, err := store.BatchSelectEmployeesByEmails([]string{email})
	if err != nil {
		return nil, err
	}

	if len(employees) == 0 {
		return nil, store.DataNotFoundError{}
	}

	employee := employees[0]
	user := mapEmployeeEntity2User(employee)

	return &user, nil
}

func (builder UserCmd) BuildEmployees() ([]User, error) {
	result := make([]User, 0)

	employees, err := store.SelectEmployees(builder.Options)
	if err != nil {
		return nil, err
	}

	for _, employee := range *employees {
		employee := employee
		user := mapEmployeeEntity2User(employee)
		result = append(result, user)
	}

	return result, nil
}

func (builder UserCmd) BuildUsersByGhLogins(logins []string) (map[string]User, error) {
	result := make(map[string]User)
	employees, err := store.BatchSelectEmployeesByGhLogins(logins)
	if err != nil {
		return nil, err
	}

	for _, employee := range employees {
		employee := employee
		user := mapEmployeeEntity2User(employee)
		result[user.GitLogin] = user
	}

	for _, githubId := range logins {
		if _, ok := result[githubId]; !ok {
			result[githubId] = User{
				IsEmployee: false,
				GitUser: GitUser{
					GitLogin: githubId,
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
