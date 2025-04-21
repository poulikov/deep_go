package main

import (
	"reflect"
	"runtime"
	"sync"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type COWBuffer struct {
	data []byte
	refs *int

	mx *sync.RWMutex
}

func NewCOWBuffer(data []byte) COWBuffer {
	buf := COWBuffer{
		data: data,
		refs: new(int),
		mx:   &sync.RWMutex{},
	}
	c := clean{
		refs: buf.refs,
		mx:   buf.mx,
	}
	runtime.AddCleanup(&buf, func(c clean) {
		c.mx.Lock()
		if *c.refs > 0 {
			*c.refs--
		}
		c.mx.Unlock()
	}, c)
	return buf
}

func (b *COWBuffer) Clone() COWBuffer {
	b.mx.Lock()
	defer b.mx.Unlock()
	*b.refs++
	buf := COWBuffer{
		data: b.data,
		refs: b.refs,
		mx:   b.mx,
	}
	c := clean{
		refs: buf.refs,
		mx:   buf.mx,
	}
	runtime.AddCleanup(&buf, func(c clean) {
		c.mx.Lock()
		if *c.refs > 0 {
			*c.refs--
		}
		c.mx.Unlock()
	}, c)
	return buf
}

func (b *COWBuffer) Close() {
	b.mx.Lock()
	if *b.refs == 0 {
		b.mx.Unlock()
		return
	}
	*b.refs--
	b.refs = new(int)
	b.data = nil
	defer b.mx.Unlock()
	b.mx = &sync.RWMutex{}
}

func (b *COWBuffer) Update(index int, value byte) bool {
	b.mx.Lock()
	unlock := b.mx.Unlock

	if index < 0 || index >= len(b.data) {
		unlock()
		return false
	}
	if *b.refs > 0 {
		*b.refs--
		data := make([]byte, len(b.data))
		copy(data, b.data)
		b.data = data
		b.refs = new(int)
		mx := &sync.RWMutex{}
		mx.Lock()
		defer mx.Unlock()
		b.mx = mx
		c := clean{
			refs: b.refs,
			mx:   b.mx,
		}
		runtime.AddCleanup(&b, func(c clean) {
			c.mx.Lock()
			if *c.refs > 0 {
				*c.refs--
			}
			c.mx.Unlock()
		}, c)
	}
	b.data[index] = value
	unlock()

	return true
}

func (b *COWBuffer) String() string {
	b.mx.RLock()
	defer b.mx.RUnlock()
	if len(b.data) == 0 {
		return ""
	}
	return unsafe.String(unsafe.SliceData(b.data), len(b.data))
}

type clean struct {
	refs *int
	mx   *sync.RWMutex
}

func TestCOWBuffer(t *testing.T) {
	data := []byte{'a', 'b', 'c', 'd'}
	buffer := NewCOWBuffer(data)
	defer buffer.Close()

	copy1 := buffer.Clone()
	copy2 := buffer.Clone()

	assert.Equal(t, unsafe.SliceData(data), unsafe.SliceData(buffer.data))
	assert.Equal(t, unsafe.SliceData(buffer.data), unsafe.SliceData(copy1.data))
	assert.Equal(t, unsafe.SliceData(copy1.data), unsafe.SliceData(copy2.data))

	assert.True(t, (*byte)(unsafe.SliceData(data)) == unsafe.StringData(buffer.String()))
	assert.True(t, (*byte)(unsafe.StringData(buffer.String())) == unsafe.StringData(copy1.String()))
	assert.True(t, (*byte)(unsafe.StringData(copy1.String())) == unsafe.StringData(copy2.String()))

	assert.True(t, buffer.Update(0, 'g'))
	assert.False(t, buffer.Update(-1, 'g'))
	assert.False(t, buffer.Update(4, 'g'))

	assert.True(t, reflect.DeepEqual([]byte{'g', 'b', 'c', 'd'}, buffer.data))
	assert.True(t, reflect.DeepEqual([]byte{'a', 'b', 'c', 'd'}, copy1.data))
	assert.True(t, reflect.DeepEqual([]byte{'a', 'b', 'c', 'd'}, copy2.data))

	assert.NotEqual(t, unsafe.SliceData(buffer.data), unsafe.SliceData(copy1.data))
	assert.Equal(t, unsafe.SliceData(copy1.data), unsafe.SliceData(copy2.data))

	copy1.Close()

	previous := copy2.data
	copy2.Update(0, 'f')
	current := copy2.data

	// 1 reference - don't need to copy buffer during update
	assert.Equal(t, unsafe.SliceData(previous), unsafe.SliceData(current))

	copy2.Close()
}
