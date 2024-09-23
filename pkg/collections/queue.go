package collections

import "fmt"

type Queue[T any] struct {
	data []T
}

func NewQueue[T any]() *Queue[T] {
	q := Queue[T]{}
	q.data = make([]T, 0)
	return &q
}

func (q *Queue[T]) Queue(item T) {
	q.data = append([]T{item}, q.data...)
}

func (q *Queue[T]) Print() {
	fmt.Printf("q.data: %v\n", q.data)
}

func (q *Queue[T]) Dequeue() T {
	first := q.data[0]
	if len(q.data) >= 2 {
		q.data = q.data[1:len(q.data)]
	} else {
		q.data = q.data[0:len(q.data)]
	}
	return first
}

func (q *Queue[T]) Peek() T {
	first := q.data[len(q.data)-1]
	return first
}

func (q *Queue[_]) Empty() bool {
	return len(q.data) == 0
}

func (q *Queue[_]) Len() int {
	return len(q.data)
}
