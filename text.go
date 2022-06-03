package hype

import "fmt"

type Text string

func (tn Text) Children() Nodes {
	return Nodes{}
}

func (tn Text) String() string {
	return string(tn)
}

func (tn Text) Format(f fmt.State, verb rune) {
	if len(tn) == 0 {
		return
	}

	switch verb {
	case 'v':
	default:
		fmt.Fprintf(f, "%s", tn.String())
	}
}
