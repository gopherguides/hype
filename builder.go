package hype

// func elementBuilder(node *Node) string {
// 	atom := node.DataAtom.String()

// 	sb := strings.Builder{}
// 	sb.WriteString("<" + atom)

// 	ats := node.attrs.String()
// 	if len(ats) > 0 {
// 		sb.WriteString(" " + ats)
// 	}
// 	sb.WriteString(">")

// 	cs := node.Children.String()
// 	if len(cs) > 0 {
// 		sb.WriteString(cs)
// 	}

// 	sb.WriteString("</" + atom + ">")
// 	return sb.String()
// }

// func inlineBuilder(node *Node) string {
// 	atom := node.DataAtom.String()

// 	sb := strings.Builder{}
// 	sb.WriteString("<" + atom)

// 	ats := node.attrs.String()
// 	if len(ats) > 0 {
// 		sb.WriteString(" " + ats)
// 	}
// 	sb.WriteString(" />")
// 	return sb.String()
// }
