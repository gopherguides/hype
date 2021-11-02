package main

type Venue struct{}

// snippet: entertainer
type Entertainer interface {
	// snippet: entertainer-funcs
	Name() string
	Perform(v Venue) error
	// snippet: entertainer-funcs
}

// snippet: entertainer

type Setuper interface {
	// snippet: setuper-funcs
	Setup(v Venue) error
	// snippet: setuper-funcs
}

type Teardowner interface {
	// snippet: teardowner-funcs
	Teardown(v Venue) error
	// snippet: teardowner-funcs
}
