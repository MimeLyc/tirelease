package notify

import "tirelease/commons/feishu"

type NotifyContent struct {
	Header string
	// Receiver email
	Blocks []Block
}

type Block struct {
	Text  string
	Links []Link
}

type Link struct {
	Href string
	Text string
}

func (content NotifyContent) ParseToFeishuContent() [][]feishu.ContentElement {
	result := make([][]feishu.ContentElement, 0)
	blocks := content.Blocks
	for _, block := range blocks {
		result = append(result,
			[]feishu.ContentElement{feishu.NewTextContentElement(block.Text)},
		)
		for _, link := range block.Links {
			result = append(result,
				[]feishu.ContentElement{feishu.NewHrefContentElement(link.Href, link.Text)},
			)
		}
	}
	return result
}
