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

func (s *Stack[T]) Pop() T {
	last := s.data[len(s.data)-1]
	s.data = s.data[0 : len(s.data)-1]
	return last
}

func (s *Stack[T]) Peek() T {
	last := s.data[len(s.data)-1]
	return last
}

func (s *Stack[_]) Empty() bool {
	return len(s.data) == 0
}

func (s *Stack[_]) Len() int {
	return len(s.data)
}
