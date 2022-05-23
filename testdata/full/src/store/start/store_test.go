package demo

import "testing"

// snippet: example
func Test_Store_Insert(t *testing.T) {
	t.Parallel()

	// create a store
	s := &Store{
		data: Data{},
	}

	exp := 1

	// insert a non-valid type
	err := s.Insert(exp, func() {})
	if err != nil {
		t.Fatal(err)
	}

	// retreive the type
	act, err := s.Find(exp)
	if err != nil {
		t.Fatal(err)
	}

	// assert the returned value is a func()
	_, ok := act.(func())
	if !ok {
		t.Fatalf("unexpected type %T", act)
	}

}

// snippet: example
