package main

import (
	"cmp"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type kvpair[K cmp.Ordered, V any] struct {
	key   K
	value V
	left  *kvpair[K, V]
	right *kvpair[K, V]
}

type OrderedMap[K cmp.Ordered, V any] struct {
	// need to implement
	root *kvpair[K, V]
	size int
}

func NewOrderedMap[K cmp.Ordered, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{} // need to implement
}

func (m *OrderedMap[K, V]) Insert(key K, value V) {
	newPair := &kvpair[K, V]{key: key, value: value}

	if m.root == nil {
		m.root = newPair
		m.size++
		return
	}

	current := m.root
	for current != nil {
		switch {
		case key == current.key:
			current.value = value
			return
		case key < current.key:
			if current.left == nil {
				current.left = newPair
				m.size++
				return
			}
			current = current.left
		case key > current.key:
			if current.right == nil {
				current.right = newPair
				m.size++
				return
			}
			current = current.right
		}
	}
}

func (m *OrderedMap[K, V]) Erase(key K) {
	var parent, current *kvpair[K, V]
	current = m.root

	for current != nil && current.key != key {
		parent = current
		if current.key > key {
			current = current.left
		} else {
			current = current.right
		}
	}
	if current == nil {
		return
	}
	if parent == nil {
		if current.key == key {
			m.root = nil
		}
		return
	}

	switch {
	case current.left == nil && current.right == nil:
		if parent.key > current.key {
			parent.left = nil
		} else {
			parent.right = nil
		}
	case current.left == nil:
		*current = *current.right
	case current.right == nil:
		*current = *current.left
	default:
		pair := m.findMin(current.right)
		m.Erase(pair.key)
		current.value = pair.value
		current.key = pair.key
	}

	m.size--
}

func (m *OrderedMap[K, V]) findMin(pair *kvpair[K, V]) *kvpair[K, V] {
	current := pair
	for current.left != nil {
		current = current.left
	}
	return current
}

func (m *OrderedMap[K, V]) Contains(key K) bool {
	current := m.root

	for current != nil {
		switch {
		case key == current.key:
			return true
		case key < current.key:
			current = current.left
		case key > current.key:
			current = current.right
		}
	}

	return false
}

func (m *OrderedMap[K, V]) Size() int {
	return m.size
}

func (m *OrderedMap[K, V]) ForEach(action func(K, V)) {
	m.walk(m.root, action)
}

func (m *OrderedMap[K, V]) walk(pair *kvpair[K, V], f func(K, V)) {
	if pair != nil {
		m.walk(pair.left, f)
		f(pair.key, pair.value)
		m.walk(pair.right, f)
	}
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap[int, int]()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
