// Tool Url: https://github.com/jinzhu/configor

package configs

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path"
)

const (
	TestConfig       = "tests/test-config/config.yaml"
	TestSecretConfig = "tests/test-config/secrets"
)

// Config Database configuration
type Config struct {
	Secret
}

// NewConfig Load config from file into 'Config' variable
func NewConfig(configPath, secretDir string) *Config {
	content, err := os.ReadFile(configPath)
	if err != nil {
		panic(fmt.Sprintf("read config file failed, err= %v", err))
	}

	c := &Config{}
	if err := yaml.Unmarshal(content, c); err != nil {
		panic(fmt.Sprintf("parse config file failed, err= %v", err))
	}

	// read secret config
	readSecret := func(secret string) []byte {
		content, err = os.ReadFile(path.Join(secretDir, secret))
		if err != nil {
			panic(fmt.Sprintf("read secret config %s failed, err= %v", secret, err))
		}
		return content
	}
	c.DSN = string(readSecret("dsn"))
	c.EmployeeDSN = string(readSecret("employeeDSN"))
	c.GitHubAccessToken = string(readSecret("gitHubAccessToken"))
	if err = yaml.Unmarshal(readSecret("feiShu"), &c.FeiShu); err != nil {
		panic(fmt.Sprintf("read feiShu config failed, err= %v", err))
	}

	return c
}

type Secret struct {
	DSN               string `yaml:"dsn"`
	EmployeeDSN       string `yaml:"employeeDSN"`
	GitHubAccessToken string `yaml:"gitHubAccessToken"`
	FeiShu            struct {
		AppId     string `yaml:"appId"`
		AppSecret string `yaml:"appSecret"`
	} `yaml:"feiShu"`
}
