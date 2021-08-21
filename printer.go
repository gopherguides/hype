package hype

import "io"

func Print(tags Tags, w io.Writer, fn func(io.Writer, Tag) error) error {
	for _, tag := range tags {
		if err := fn(w, tag); err != nil {
			return err
		}
		if err := Print(tag.GetChildren(), w, fn); err != nil {
			return err
		}
	}
	return nil
}
