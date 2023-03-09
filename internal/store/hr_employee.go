package store

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
	"tirelease/internal/entity"
)

// ATTENTION: there must not be any create/update/insert operation in this file.
var DB *gorm.DB

func SelectAllHrEmployee() ([]entity.HrEmployee, error) {
	var hrEmployees []entity.HrEmployee
	result := DB.Find(&hrEmployees)
	if result.Error != nil {
		return nil, result.Error
	}

	return hrEmployees, nil
}

func InitHrEmployeeDB(dsn string) {
	// Connect
	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err.Error())
	}
	sqlDB, err := conn.DB()
	if err != nil {
		panic(err.Error())
	}
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Second * 600)

	DB = conn
}
