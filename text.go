package hype

type TextNode interface {
	Node
	Text() string
}

var _ TextNode = Text("")

type Text string

func (tn Text) Children() Nodes {
	return Nodes{}
}

func (tn Text) String() string {
	return string(tn)
}

func (tn Text) Text() string {
	return string(tn)
}
