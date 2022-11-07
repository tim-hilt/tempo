package set

import "fmt"

type set[T comparable] struct {
	items map[T]bool
}

func New[T comparable]() set[T] {
	return set[T]{items: make(map[T]bool)}
}

func (s *set[T]) Delete(item T) {
	delete(s.items, item)
}

func (s *set[T]) Add(item T) {
	s.items[item] = true
}

func (s set[T]) Contains(item T) bool {
	_, ok := s.items[item]
	return ok
}

func (s *set[T]) Reset() {
	s.items = make(map[T]bool)
}

func (s set[T]) Items() []T {
	items := make([]T, len(s.items))

	i := 0
	for k := range s.items {
		items[i] = k
		i++
	}
	return items
}

func (s set[T]) String() string {
	str := "{"
	for item := range s.items {
		str = str + fmt.Sprint(item) + ", "
	}
	str = str[:len(str)-2] + "}"
	return str
}
