package collections_test

import (
	"graph/pkg/collections"
	"slices"
	"testing"
)

func TestStack(t *testing.T) {
	s := collections.NewStack[int]()

	// pop from empty stack
	for i := 0; i < 3; i++ {
		n, ok := s.Pop()
		if ok {
			t.Fatalf("s.Pop() = %v, %v, want %v, %v", n, ok, 0, false)
		}
	}

	// test lifo behaviour
	numbers := []int{13, 23, 22, 9}
	for _, number := range numbers {
		s.Push(number)
	}
	slices.Reverse(numbers)
	for _, expectedNumber := range numbers {
		actualNumber, _ := s.Pop()
		if expectedNumber != actualNumber {
			t.Fatalf("s.Pop() = %v, want %v", actualNumber, expectedNumber)
		}
	}

	// push random values to the stack to test peek
	s.Push(81)
	s.Push(2)

	// test peek behaviour
	expectedNumber := 10
	s.Push(expectedNumber)
	actualNumber, ok := s.Peek()
	if expectedNumber != actualNumber {
		t.Fatalf("s.Pop() = %v, %v, want %v, %v", actualNumber, ok, expectedNumber, true)
	}

	// test final length
	queueSize := s.Len()
	expectedQueueSize := 3
	if queueSize != expectedQueueSize {
		t.Fatalf("s.Len() = %v, want %v", queueSize, expectedQueueSize)
	}
}
