package set

type set[T comparable] struct {
	items map[T]bool
}

func New[T comparable]() set[T] {
	s := set[T]{}
	s.items = make(map[T]bool)
	return s
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
