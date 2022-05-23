package demo

import "fmt"

type Data map[int]Model

type Store struct {
	data Data
}

// snippet: func
func (s *Store) Insert(m Model) error {
	// snippet: func
	if s.data == nil {
		s.data = Data{}
	}

	s.data[m.ID()] = m
	return nil
}

func (s Store) Find(id int) (Model, error) {
	if s.data == nil {
		return nil, fmt.Errorf("key not found %d", id)
	}

	m, ok := s.data[id]
	if !ok {
		return nil, fmt.Errorf("key not found %d", id)
	}

	return m, nil
}
