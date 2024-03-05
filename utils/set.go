package utils

// Convenience wrapper around map[T]struct{} for more easily working with sets of values
type Set[T comparable] struct {
	underlying map[T]struct{}
}

func (s *Set[T]) ensureInitialized() {
	if s == nil {
		s = &Set[T]{}
	}

	if s.underlying == nil {
		s.underlying = make(map[T]struct{})
	}
}

// non-mutating methods

func (s *Set[T]) Has(element T) bool {
	if s == nil {
		return false
	}

	_, exists := s.underlying[element]
	return exists
}

func (s *Set[T]) Size() int {
	if s == nil || s.underlying == nil {
		return 0
	}

	return len(s.underlying)
}

func (s *Set[T]) Elements() []T {
	elements := []T{}

	if s != nil && s.underlying != nil {
		for element := range s.underlying {
			elements = append(elements, element)
		}
	}

	return elements
}

// mutating methods

func (s *Set[T]) Add(element T) {
	s.ensureInitialized()

	s.underlying[element] = struct{}{}
}

// removes `element` from the set
// returns true iff `element` was in the set beforehand, false if it wasn't
func (s *Set[T]) Delete(element T) bool {
	if s == nil {
		return false
	}

	elemWasInSet := s.Has(element)

	delete(s.underlying, element)
	return elemWasInSet
}

// TODO - Clear() method if I need it, that removes all elements?
