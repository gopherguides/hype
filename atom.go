package hype

type Adam string

func (a Adam) String() string {
	return string(a)
}

func (a Adam) IsZero() bool {
	return len(a) == 0
}

const (
	ERROR_ADAM      Adam = "ERROR"
	Code_Adam       Adam = "code"
	File_Adam       Adam = "file"
	File_Group_Adam Adam = "filegroup"
	Include_Adam    Adam = "include"
	Page_Adam       Adam = "page"
)

type Adamable interface {
	Adam() Adam
}

func IsAdam(a Adamable, wants ...Adam) bool {
	if a == nil {
		return false
	}

	for _, want := range wants {
		if a.Adam() == want {
			return true
		}
	}
	return false
}
