package feishu

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestMsgContentMarshal(t *testing.T) {
	contentPost := ContentPost{
		ZhCnContent: ContentWrapper{
			Title: "Test",
			Rows: [][]ContentElement{
				{
					NewTextContentElement("test11"),
					NewHrefContentElement("test12", "test12"),
				},
				{

					NewTextContentElement("test21"),
				},
			},
		},
	}

	jsonContent, _ := json.Marshal(contentPost)
	fmt.Print(jsonContent)
}
