package demo

// snippet: callbacks
type BeforeInsertable interface {
	BeforeInsert() error
}

type AfterInsertable interface {
	AfterInsert() error
}

// snippet: callbacks
