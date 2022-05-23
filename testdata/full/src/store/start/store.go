package demo

import "fmt"

type Data map[int]any

type Store struct {
	data Data
}

// snippet: func
func (s *Store) Insert(id int, m any) error {
	// snippet: func

	s.data[id] = m
	return nil
}

func (s Store) Find(id int) (any, error) {
	if s.data == nil {
		return nil, fmt.Errorf("key not found %d", id)
	}

	m, ok := s.data[id]
	if !ok {
		return nil, fmt.Errorf("key not found %d", id)
	}

	return m, nil
}
