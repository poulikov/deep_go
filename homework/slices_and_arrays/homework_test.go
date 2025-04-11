package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type CircularQueue[T ~int8 | ~int16 | ~int32 | ~int64 | ~int] struct {
	values   []T
	head     int
	tail     int
	count    int
	capacity int
}

func NewCircularQueue[T ~int8 | ~int16 | ~int32 | ~int64 | ~int](size int) *CircularQueue[T] {
	return &CircularQueue[T]{
		values:   make([]T, size),
		capacity: size,
	}
}

func (q *CircularQueue[T]) Push(value T) bool {
	if q.count == q.capacity {
		return false
	}
	q.values[q.tail] = value
	q.tail = (q.tail + 1) % q.capacity
	q.count++
	return true
}

func (q *CircularQueue[T]) Pop() bool {
	if q.count == 0 {
		return false
	}
	q.head = (q.head + 1) % q.capacity
	q.count--
	return true
}

func (q *CircularQueue[T]) Front() T {
	if q.count == 0 {
		return -1
	}
	return q.values[q.head]
}

func (q *CircularQueue[T]) Back() T {
	if q.count == 0 {
		return -1
	}
	return q.values[(q.tail-1+q.capacity)%q.capacity]
}

func (q *CircularQueue[T]) Empty() bool {
	return q.count == 0
}

func (q *CircularQueue[T]) Full() bool {
	return q.count == q.capacity
}

func TestCircularQueue(t *testing.T) {
	const queueSize = 3
	queue := NewCircularQueue[int](queueSize)

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.Equal(t, -1, queue.Front())
	assert.Equal(t, -1, queue.Back())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Push(1))
	assert.True(t, queue.Push(2))
	assert.True(t, queue.Push(3))
	assert.False(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{1, 2, 3}, queue.values))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	assert.Equal(t, 1, queue.Front())
	assert.Equal(t, 3, queue.Back())

	assert.True(t, queue.Pop())
	assert.False(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.True(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{4, 2, 3}, queue.values))

	assert.Equal(t, 2, queue.Front())
	assert.Equal(t, 4, queue.Back())

	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())
}
