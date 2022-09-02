package service

type NotifyContentTemplate struct {
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
