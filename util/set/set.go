package set

import "fmt"

type Set[T comparable] struct {
	items map[T]bool
}

func New[T comparable](is ...T) Set[T] {
	s := Set[T]{items: make(map[T]bool)}
	for _, i := range is {
		s.Add(i)
	}
	return s
}

func (s *Set[T]) Delete(item T) {
	delete(s.items, item)
}

func (s *Set[T]) Add(item T) {
	s.items[item] = true
}

func (s Set[T]) Contains(item T) bool {
	_, ok := s.items[item]
	return ok
}

func (s *Set[T]) Reset() {
	s.items = make(map[T]bool)
}

func (s Set[T]) Items() []T {
	items := make([]T, len(s.items))

	i := 0
	for k := range s.items {
		items[i] = k
		i++
	}
	return items
}

func (s Set[T]) String() string {
	str := "{"
	for item := range s.items {
		str = str + fmt.Sprint(item) + ", "
	}
	str = str[:len(str)-2] + "}"
	return str
}

func (s Set[T]) Cardinality() int {
	return len(s.items)
}

func (s1 Set[T]) Union(s2 Set[T]) Set[T] {
	newSet := New[T]()
	for item := range s1.items {
		newSet.Add(item)
	}
	for item := range s2.items {
		newSet.Add(item)
	}
	return newSet
}

func (s1 Set[T]) Intersect(s2 Set[T]) Set[T] {
	newSet := New[T]()
	for item := range s1.items {
		if s2.Contains(item) {
			newSet.Add(item)
		}
	}
	return newSet
}
