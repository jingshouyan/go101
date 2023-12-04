package util

type Set[T comparable] map[T]struct{}

func NewSet[T comparable]() Set[T] {
	return make(Set[T])
}

func (s Set[T]) Add(t T) {
	s[t] = struct{}{}
}

func (s Set[T]) Remove(t T) {
	delete(s, t)
}

func (s Set[T]) Has(t T) bool {
	_, ok := s[t]
	return ok
}

func (s Set[T]) Len() int {
	return len(s)
}

func (s Set[T]) Values() []T {
	values := make([]T, 0, len(s))
	for v := range s {
		values = append(values, v)
	}
	return values
}

func (s Set[T]) Clear() {
	for k := range s {
		delete(s, k)
	}
}

func (s Set[T]) Copy() Set[T] {
	newSet := NewSet[T]()
	for k := range s {
		newSet.Add(k)
	}
	return newSet
}

func (s Set[T]) Equal(other Set[T]) bool {
	if s.Len() != other.Len() {
		return false
	}
	for k := range s {
		if !other.Has(k) {
			return false
		}
	}
	return true
}

func (s Set[T]) Intersect(other Set[T]) Set[T] {
	newSet := NewSet[T]()
	for k := range s {
		if other.Has(k) {
			newSet.Add(k)
		}
	}
	return newSet
}

func (s Set[T]) Union(other Set[T]) Set[T] {
	newSet := NewSet[T]()
	for k := range s {
		newSet.Add(k)
	}
	for k := range other {
		newSet.Add(k)
	}
	return newSet
}

func (s Set[T]) Difference(other Set[T]) Set[T] {
	newSet := NewSet[T]()
	for k := range s {
		if !other.Has(k) {
			newSet.Add(k)
		}
	}
	return newSet
}
