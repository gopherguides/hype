package hype

type Snippet struct {
	Content  string
	File     string
	Language string
	Name     string
	Start    int
	End      int
}

func (snip Snippet) String() string {
	return snip.Content
}

type Snippets map[string]Snippet
