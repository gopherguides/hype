package hype

type TextOnlyNode interface {
	Node
	Text() string
}

var _ TextOnlyNode = TextNode("")

type TextNode string

func (tn TextNode) Children() Nodes {
	return Nodes{}
}

func (tn TextNode) String() string {
	return string(tn)
}

func (tn TextNode) Text() string {
	return string(tn)
}
