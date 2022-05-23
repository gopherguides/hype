package demo

import "fmt"

type Data map[int]Model

type Store struct {
	data Data
}

// snippet: insert
// snippet: func
func (s *Store) Insert(m Validatable) error {
	// snippet: func

	// before insert

	// validate model
	err := m.Validate()
	if err != nil {
		return err
	}

	// insert
	s.data[m.ID()] = m

	// after insert

	return nil
}

// snippet: insert

func (s Store) Find(id int) (Model, error) {
	m, ok := s.data[id]
	if !ok {
		return nil, fmt.Errorf("key not found %d", id)
	}

	return m, nil
}
