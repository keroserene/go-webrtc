package webrtc

import "sync"

// A general map that stores pointers by an unique int index, since we can't
// keep Go pointers in C land.  This is a little workaround suggested in the
// proposal for the cgo rules as of 1.6

type CGOMap struct {
	index    int
	pointers map[int]interface{}
	lock     *sync.Mutex
}

func NewCGOMap() CGOMap {
	var m CGOMap
	m.pointers = make(map[int]interface{})
	m.lock = &sync.Mutex{}
	return m
}

func (m *CGOMap) Get(index int) interface{} {
	p, ok := m.pointers[index]
	if p == nil || !ok {
		panic("couldn't retrieve the pointer")
	}
	return p
}

func (m *CGOMap) Set(p interface{}) int {
	m.lock.Lock()
	m.index += 1
	m.pointers[m.index] = p
	m.lock.Unlock()
	return m.index
}

func (m *CGOMap) Delete(index int) {
	m.lock.Lock()
	delete(m.pointers, index)
	m.lock.Unlock()
}
