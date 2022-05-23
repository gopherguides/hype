package demo

import "fmt"

type Data map[int]Model

type Store struct {
	data Data
}

func (s *Store) Insert(m Validatable) error {

	// snippet: before
	if bi, ok := m.(BeforeInsertable); ok {
		if err := bi.BeforeInsert(); err != nil {
			return err
		}
	}
	// snippet: before

	// validate model
	err := m.Validate()
	if err != nil {
		return err
	}

	// insert
	s.data[m.ID()] = m

	// after insert

	// snippet: after
	if ai, ok := m.(AfterInsertable); ok {
		if err := ai.AfterInsert(); err != nil {
			return err
		}
	}
	// snippet: after

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
