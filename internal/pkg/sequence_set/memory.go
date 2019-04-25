package set

func NewMemorySet() Set {
	return &inMemorySet{
		add:    map[string]int{},
		remove: map[string]int{},
	}
}

type inMemorySet struct {
	add    map[string]int
	remove map[string]int
}

func (s *inMemorySet) Add(key string, sequenceNumber int) error {
	a, added := s.add[key]
	if !added || a < sequenceNumber {
		s.add[key] = sequenceNumber
	}

	return nil
}

func (s *inMemorySet) Remove(key string, sequenceNumber int) error {
	r, removed := s.remove[key]
	if !removed || r < sequenceNumber {
		s.remove[key] = sequenceNumber
	}

	return nil
}

func (s *inMemorySet) Has(key string) (bool, error) {
	a := s.add[key]
	r := s.remove[key]

	return a > r, nil
}
