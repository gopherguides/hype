package assignment04

type Entertainer interface {
	Name() string
	Perform(v Venue) error
}

type Setuper interface {
	Setup(v Venue) error
}

type Teardowner interface {
	Teardown(v Venue) error
}
