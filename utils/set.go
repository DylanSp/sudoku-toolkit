package utils

import "maps"

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

// shallow clone
func (s *Set[T]) Clone() Set[T] {
	newSet := Set[T]{}

	newSet.underlying = maps.Clone(s.underlying)

	return newSet
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

	s.ensureInitialized()

	elemWasInSet := s.Has(element)

	delete(s.underlying, element)
	return elemWasInSet
}

// removes all elements from the set
func (s *Set[T]) DeleteAll() {
	s.ensureInitialized()

	s.underlying = make(map[T]struct{})
}
