package demo

import (
	"testing"
)

// snippet: example
func Test_Store_Insert(t *testing.T) {
	t.Parallel()

	// create a store
	s := &Store[string, User]{
		data: map[string]User{},
	}

	// create a user
	exp := User{Email: "kurt@exampl.com"}

	// insert the user
	err := s.Insert(exp)
	if err != nil {
		t.Fatal(err)
	}

	// retreive the user
	act, err := s.Find(exp.Email)
	if err != nil {
		t.Fatal(err)
	}

	// assert the returned user is the same as the inserted user
	if exp.Email != act.Email {
		t.Fatalf("expected %v, got %v", exp, act)
	}

}

// snippet: example
