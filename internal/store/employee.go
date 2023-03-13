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

func SelectEmployees(options *entity.EmployeeOptions) (*[]entity.Employee, error) {
	sql := "select * from employees where 1=1" + employeeWhere(options) + employeeOrderBy(options) + employeeLimit(options)
	// 查询
	var employees []entity.Employee
	if err := storeGlobalDB.RawWrapper(sql, options).Find(&employees).Error; err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("find employees: %+v failed", options))
	}

	return &employees, nil
}

func employeeWhere(option *entity.EmployeeOptions) string {
	sql := ""
	if option.HrEmployeeID != "" {
		sql += " and employees.hr_employee_id = @HrEmployeeID"
	}
	if option.Name != "" {
		sql += " and employees.name = @Name"
	}
	if option.Email != "" {
		sql += " and employees.email = @Email"
	}
	if option.GithubId != "" {
		sql += " and employees.github_id = @GithubId"
	}
	if option.IsActive != nil {
		sql += " and employees.active = @IsActive"
	}
	if option.GhEmail != "" {
		sql += " and employees.gh_email = @GhEmail"
	}
	if option.GhName != "" {
		sql += " and employees.gh_name = @GhName"
	}
	if option.GithubIds != nil {
		sql += " and employees.github_id in @GithubIds"
	}
	if option.Emails != nil {
		sql += " and employees.email in @Emails"
	}
	return sql
}

func employeeOrderBy(option *entity.EmployeeOptions) string {
	return option.GetOrderByString()
}

func employeeLimit(option *entity.EmployeeOptions) string {
	return option.GetLimitString()
}
