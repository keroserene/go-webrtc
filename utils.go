package webrtc

import (
	"strconv"
	"sync"
)

// CGOMap is a general map that stores pointers by an unique int index, since we can't
// keep Go pointers in C land.  This is a little workaround suggested in the
// proposal for the cgo rules as of 1.6
type CGOMap struct {
	index    int
	pointers map[int]interface{}
	lock     *sync.Mutex
}

// NewCGOMap construct map
func NewCGOMap() CGOMap {
	var m CGOMap
	m.pointers = make(map[int]interface{})
	m.lock = &sync.Mutex{}
	return m
}

// Get returns value of one key
func (m *CGOMap) Get(index int) interface{} {
	p, ok := m.pointers[index]
	if p == nil || !ok {
		panic("couldn't retrieve the pointer")
	}
	return p
}

// Set set key value pair
func (m *CGOMap) Set(p interface{}) int {
	m.lock.Lock()
	m.index++
	m.pointers[m.index] = p
	m.lock.Unlock()
	return m.index
}

// Delete delete one key value pair
func (m *CGOMap) Delete(index int) {
	m.lock.Lock()
	delete(m.pointers, index)
	m.lock.Unlock()
}

// EnumToStringSafe return a string value for an integer enum from a mapping array
// or the integer string if the integer it outside the expected range.
func EnumToStringSafe(value int, valueStrings []string) string {
	if value >= 0 && value < len(valueStrings) {
		return valueStrings[value]
	}
	return strconv.Itoa(value)
}
