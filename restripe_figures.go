package hype

type IDGenerator func(i int, fig *Figure) (string, error)

// RestripeFigureIDs will rewrite all of the figure IDs,
// and they're references, using the given IDGenerator.
func RestripeFigureIDs(nodes Nodes, fn IDGenerator) error {
	if fn == nil {
		return ErrIsNil("IDGenerator")
	}

	figs := ByType[*Figure](nodes)

	for i, fig := range figs {

		fid, err := fig.ValidAttr("id")
		if err != nil {
			return err
		}

		uid, err := fn(i, fig)
		if err != nil {
			return err
		}

		if err := fig.Set("id", uid); err != nil {
			return err
		}

		refs := ByType[*Ref](nodes)
		for _, ref := range refs {
			rid, err := ref.ValidAttr("id")
			if err != nil {
				return err
			}

			if rid != fid {
				continue
			}

			ref.Figure = fig
			if err := ref.Set("id", uid); err != nil {
				return err
			}

		}

	}

	return nil
}
