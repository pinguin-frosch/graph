package collections_test

import (
	"graph/pkg/collections"
	"testing"
)

func TestQueue(t *testing.T) {
	q := collections.NewQueue[int]()

	// dequeue empty queue
	for i := 0; i < 3; i++ {
		n, ok := q.Dequeue()
		if ok {
			t.Fatalf("q.Dequeue() = %v, %v, want %v, %v", n, ok, 0, false)
		}
	}

	// test fifo behaviour
	numbers := []int{9, 1, 4, 7}
	for _, number := range numbers {
		q.Enqueue(number)
	}
	for _, expectedNumber := range numbers {
		actualNumber, _ := q.Dequeue()
		if expectedNumber != actualNumber {
			t.Fatalf("q.Dequeue() = %v, want %v", actualNumber, expectedNumber)
		}
	}

	// test peek behaviour
	expectedNumber := 10
	q.Enqueue(expectedNumber)

	// queue random values to test peek
	q.Enqueue(28)
	q.Enqueue(92)

	actualNumber, ok := q.Peek()
	if expectedNumber != actualNumber {
		t.Fatalf("q.Dequeue() = %v, %v, want %v, %v", actualNumber, ok, expectedNumber, true)
	}

	queueSize := q.Len()
	expectedQueueSize := 3
	if queueSize != expectedQueueSize {
		t.Fatalf("q.Len() = %v, want %v", queueSize, expectedQueueSize)
	}
}
