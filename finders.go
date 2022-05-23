package hype

func ByType[T Node](nodes Nodes) []T {
	var res []T

	for _, t := range nodes {
		if x, ok := t.(T); ok {
			res = append(res, x)
		}

		res = append(res, ByType[T](t.Children())...)
	}

	return res
}

func ByAttrs(nodes Nodes, query map[string]string) []AttrNode {
	var res []AttrNode
	for _, n := range nodes {
		t, ok := n.(AttrNode)
		if ok {
			ta := t.Attrs()

			if AttrMatches(ta, query) {
				res = append(res, t)
			}
		}

		res = append(res, ByAttrs(n.Children(), query)...)
	}
	return res
}

func ByAtom[T ~string](nodes Nodes, want ...T) []AtomableNode {
	var res []AtomableNode
	for _, n := range nodes {
		an, ok := n.(AtomableNode)
		if !ok {
			res = append(res, ByAtom(n.Children(), want...)...)
			continue
		}

		for _, w := range want {
			if an.Atom().String() == string(w) {
				res = append(res, an)
				break
			}
		}

		res = append(res, ByAtom(n.Children(), want...)...)
	}
	return res
}
