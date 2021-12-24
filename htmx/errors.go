package htmx

import "fmt"

// ErrAttrNotFound is returned when an attribute is not found.
type ErrAttrNotFound string

func (e ErrAttrNotFound) Error() string {
	return fmt.Sprintf("attribute not found: %q", string(e))
}
