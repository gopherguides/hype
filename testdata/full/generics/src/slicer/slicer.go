package demo

func Slicer[T any](input T) []T {
	return []T{input}
}
