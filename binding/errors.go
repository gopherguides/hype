package binding

import "fmt"

type ErrPath string

func (e ErrPath) Error() string {
	return fmt.Sprintf("could not parse section from: %q", string(e))
}

func (e ErrPath) Is(err error) bool {
	_, ok := err.(ErrPath)
	return ok
}
