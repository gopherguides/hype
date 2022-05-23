package demo

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

// Store is a map of models where the map key is any
// comparable type and the map value is any type that
// implements the Model constraint.
type Store[K constraints.Ordered, M Model[K]] struct {
	data map[K]M
}

// snippet: func
func (s *Store[K, M]) Insert(m M) error {
	// snippet: func

	s.data[m.ID()] = m
	return nil
}

func (s Store[K, M]) Find(id K) (M, error) {
	m, ok := s.data[id]
	if !ok {
		return m, fmt.Errorf("key not found %v", id)
	}

	return m, nil
}
