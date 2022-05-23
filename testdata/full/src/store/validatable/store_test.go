package demo

import "testing"

// snippet: example
func Test_Store_Insert(t *testing.T) {
	t.Parallel()

	// create a store
	s := &Store{
		data: Data{},
	}

	// create a user
	exp := User{UID: 1}

	// insert the user
	err := s.Insert(exp)
	if err != nil {
		t.Fatal(err)
	}

	// retreive the user
	act, err := s.Find(exp.UID)
	if err != nil {
		t.Fatal(err)
	}

	// assert the returned value is a user
	actu, ok := act.(User)
	if !ok {
		t.Fatalf("unexpected type %T", act)
	}

	// assert the returned user is the same as the inserted user
	if exp.UID != actu.UID {
		t.Fatalf("expected %v, got %v", exp, actu)
	}

}

// snippet: example
