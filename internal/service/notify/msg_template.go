package notify

import (
	"fmt"
	"tirelease/commons/feishu"
	"tirelease/internal/constants"
)

type NotifyContent struct {
	Header string
	// Receiver email
	Blocks   []Block
	Severity constants.NotifySeverity
}

type Block struct {
	Text  string
	Links []Link
}

type Link struct {
	Href string
	Text string
}

func (content NotifyContent) ParseToFeishuContent() []feishu.ContentElement {
	result := make([]feishu.ContentElement, 0)
	blocks := content.Blocks
	for _, block := range blocks {
		result = append(result,
			feishu.NewMDCardElement(block.Text),
		)
		for _, link := range block.Links {
			result = append(result,
				feishu.NewMDCardElement(
					fmt.Sprintf("<a href='%s'>%s</a>", link.Href, link.Text),
				),
			)
		}
		result = append(result,
			feishu.NewHrCardElement(),
		)
		result = append(result,
			feishu.NewFootPrintCardElement("ChatOps supported by TiRelease."),
		)
	}
	return result
}
