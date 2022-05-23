package hype

type Tag interface {
	Node
	Atomable
	StartTag() string
	EndTag() string
}

type Tags []Tag
