package assignment04

import (
	"fmt"
	"io"
)

// snippet: venue
type Venue struct {
	Audience int
	Log      io.Writer
}

// snippet: venue

func (v *Venue) Entertain(audience int, acts ...Entertainer) error {
	if len(acts) == 0 {
		return fmt.Errorf("there are no entertainers to perform")
	}

	v.Audience = audience
	for _, act := range acts {
		if err := v.play(act); err != nil {
			return err
		}
	}

	return nil
}

func (v Venue) play(act Entertainer) error {

	name := act.Name()

	if s, ok := act.(Setuper); ok {
		if err := s.Setup(v); err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
		fmt.Fprintf(v.Log, "%s has completed setup.\n", name)
	}

	if err := act.Perform(v); err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}

	fmt.Fprintf(v.Log, "%s has performed for %d people.\n", name, v.Audience)

	if t, ok := act.(Teardowner); ok {
		if err := t.Teardown(v); err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
		fmt.Fprintf(v.Log, "%s has completed teardown.\n", name)
	}

	return nil
}
