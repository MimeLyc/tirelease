package feishu_handler

import (
	"encoding/json"
	"fmt"
	"strings"
	"tirelease/internal/constants"
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

type text struct {
	Text string `json:"text,omitempty"`
}

func NewContent(raw string) (content, error) {
	text := text{}
	json.Unmarshal([]byte(raw), &text)

	contents := extractContent(text.Text)
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
		if !strings.Contains(s, "=") {
			continue
		}
		key := strings.Split(s, "=")[0]
		value := strings.Split(s, "=")[1]
		if key == "spec" {
			return value
		}
	}

	return "all"
}

func (c content) extractByKey(target string) string {
	// TODO add error detaction
	for _, s := range c.flags {
		if !strings.Contains(s, "=") {
			continue
		}
		key := strings.Split(s, "=")[0]
		value := strings.Split(s, "=")[1]
		if key == target {
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

type ActionRequest struct {
	Object   constants.EventRegisterObject `json:"register_object,omitempty"`
	ObjectID string                        `json:"register_object_id,omitempty"`
	Action   constants.EventRegisterAction `json:"register_action,omitempty"`
}

func NewActionRequest(receive ActionReceive) ActionRequest {
	object := receive.Action.Value["register_object"]
	objectId := receive.Action.Value["register_object_id"]
	action := receive.Action.Value["register_action"]

	return ActionRequest{
		Object:   constants.EventRegisterObject(object.(string)),
		ObjectID: objectId.(string),
		Action:   constants.EventRegisterAction(action.(string)),
	}
}
