// Tool Url: https://github.com/jinzhu/configor

package configs

import (
	"fmt"
	"github.com/jinzhu/configor"
)

// Database configuration
type Config struct {
	Mysql struct {
		UserName string `default:"root"`
		PassWord string `required:"true"`
		Host     string `required:"true"`
		Port     string `default:"3306"`
		DataBase string `required:"true"`
		CharSet  string `default:"utf8"`
		TimeZone string `default:"Asia%2FShanghai"`
	}

	EmployeeMysql struct {
		UserName string `default:"wh-read"`
		PassWord string `required:"true"`
		Host     string `required:"true"`
		Port     string `default:"3306"`
		DataBase string `required:"true"`
		CharSet  string `default:"utf8"`
		TimeZone string `default:"Asia%2FShanghai"`
	}

	Github struct {
		AccessToken string `required:"true"`
	}

	Feishu struct {
		AppId     string `required:"false"`
		AppSecret string `required:"false"`
	}

	Secret
}

// Load config from file into 'Config' variable
func NewConfig(file string) *Config {
	c := Config{}
	err := configor.Load(c, file)
	if err != nil {
		panic(fmt.Sprintf("parse config file failed, err= %v", err))
	}
	return &c
}

type Secret struct {
	DSN               string
	EmployeeDSN       string
	GitHubAccessToken string
}
