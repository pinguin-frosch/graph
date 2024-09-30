package collections

type Queue[T any] struct {
	data []T
}

func NewQueue[T any]() *Queue[T] {
	q := Queue[T]{}
	q.data = make([]T, 0)
	return &q
}

func (q *Queue[T]) Enqueue(item T) {
	q.data = append(q.data, item)
}

func (q *Queue[T]) Dequeue() (T, bool) {
	if q.Empty() {
		var zero T
		return zero, false
	}
	item := q.data[0]
	if len(q.data) == 1 {
		q.data = nil
	} else {
		q.data = q.data[1:]
	}
	return item, true
}

func (q *Queue[T]) Peek() (T, bool) {
	if q.Empty() {
		var zero T
		return zero, false
	}
	item := q.data[0]
	return item, true
}

func (q *Queue[T]) Empty() bool {
	return len(q.data) == 0
}

func (q *Queue[T]) Len() int {
	return len(q.data)
}
