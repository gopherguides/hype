package hype

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Text string

func (tn Text) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(map[string]any{
		"type": fmt.Sprintf("%T", tn),
		"text": string(tn),
	}, "", "  ")
}

func (tn Text) Children() Nodes {
	return Nodes{}
}

func (tn Text) String() string {
	return string(tn)
}

func (tn Text) MD() string {
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

func (tn Text) IsEmptyNode() bool {
	return len(strings.TrimSpace(string(tn))) == 0
}
