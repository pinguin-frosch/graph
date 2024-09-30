package collections

type Stack[T any] struct {
	data []T
}

func NewStack[T any]() *Stack[T] {
	s := Stack[T]{}
	s.data = make([]T, 0)
	return &s
}

func (s *Stack[T]) Push(item T) {
	s.data = append(s.data, item)
}

func (s *Stack[T]) Pop() (T, bool) {
	if s.Empty() {
		var zero T
		return zero, false
	}
	item := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return item, true
}

func (s *Stack[T]) Peek() (T, bool) {
	if s.Empty() {
		var zero T
		return zero, false
	}
	item := s.data[len(s.data)-1]
	return item, true
}

func (s *Stack[T]) Empty() bool {
	return len(s.data) == 0
}

func (s *Stack[T]) Len() int {
	return len(s.data)
}
