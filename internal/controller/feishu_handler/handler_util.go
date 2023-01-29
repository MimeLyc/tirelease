package feishu_handler

import (
	"fmt"
	"strings"
)

type Schema struct {
	Schema string `json:"schema,omitempty" form:"schema"`
}

func (s Schema) verifyVersion() bool {
	return s.Schema == "2.0"
}

type content struct {
	target string
	cmd    string
	flags  []string
}

func NewContent(raw string) (content, error) {
	contents := extractContent(raw)
	if len(contents) < 3 {
		return content{}, contentTooShortError{}
	}

	return content{
		target: contents[0],
		cmd:    contents[1],
		flags:  contents[2:],
	}, nil
}

func (c content) extractSpec() string {
	// TODO add error detaction
	for _, s := range c.flags {
		key := strings.Split(s, "=")[0]
		value := strings.Split(s, "=")[1]
		if key == "spec" {
			return value
		}
	}

	return "all"
}

func (c content) extractByKey(key string) string {
	// TODO add error detaction
	for _, s := range c.flags {
		key := strings.Split(s, "=")[0]
		value := strings.Split(s, "=")[1]
		if key == key {
			return value
		}
	}

	return ""
}

func extractContent(raw string) []string {
	formated := strings.Replace(raw, "@_user_1", "", 1)
	formated = strings.Trim(formated, " ")

	return strings.Split(formated, " ")
}

type contentTooShortError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err contentTooShortError) Error() string {
	return fmt.Sprintf(`Command is too short, the pattern is: 
        <object> <command> <options>
            object: issue, pr, version ...
            command: approve, watch ...
    `)
}
