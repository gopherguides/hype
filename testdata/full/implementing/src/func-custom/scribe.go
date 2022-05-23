package main

type Scribe struct {
	data []byte
}

func (s *Scribe) Write(p []byte) (int, error) {
	s.data = p
	return len(p), nil
}

func (s Scribe) String() string {
	return string(s.data)
}
