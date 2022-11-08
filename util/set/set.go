package set

import "fmt"

type Set[T comparable] struct {
	items map[T]bool
}

func New[T comparable]() Set[T] {
	return Set[T]{items: make(map[T]bool)}
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
