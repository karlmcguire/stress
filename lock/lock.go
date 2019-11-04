package lock

import "sync"

type Map struct {
	sync.RWMutex
	data map[uint64]uint64
}

func New(size uint64) *Map {
	return &Map{
		data: make(map[uint64]uint64, size),
	}
}

func (m *Map) Get(key uint64) uint64 {
	m.RLock()
	val := m.data[key]
	m.RUnlock()
	return val
}

func (m *Map) Set(key, val uint64) {
	m.Lock()
	m.data[key] = val
	m.Unlock()
}

func (m *Map) SetAll(pairs [][2]uint64) {
	m.Lock()
	for _, pair := range pairs {
		m.data[pair[0]] = pair[1]
	}
	m.Unlock()
}

func (m *Map) Del(key uint64) {
	m.Lock()
	delete(m.data, key)
	m.Unlock()
}
