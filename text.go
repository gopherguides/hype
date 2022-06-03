package hype

type Text string

func (tn Text) Children() Nodes {
	return Nodes{}
}

func (tn Text) String() string {
	return string(tn)
}
