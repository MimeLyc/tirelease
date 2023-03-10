// Tool Url: https://github.com/go-gorm/gorm
// Tool Guide: https://gorm.io/docs/

package database

import (
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	*gorm.DB
}

func MustConnect(dsn string) *Database {
	// Connect
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err.Error())
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err.Error())
	}
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Second * 600)

	return &Database{DB: db}
}

func (db *Database) RawWrapper(sql string, values ...interface{}) (tx *gorm.DB) {
	if strings.Contains(sql, "@") || strings.Contains(sql, "?") {
		return db.DB.Raw(sql, values...)
	} else {
		return db.DB.Raw(sql)
	}
}
