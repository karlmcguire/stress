package chanDrop

import "sync"

type Map struct {
	sync.RWMutex
	items chan [2]uint64
	data  map[uint64]uint64
}

func New(size uint64) *Map {
	m := &Map{
		items: make(chan [2]uint64, size/10),
		data:  make(map[uint64]uint64, size),
	}
	go m.process()
	return m
}

func (m *Map) Get(key uint64) uint64 {
	m.RLock()
	val := m.data[key]
	m.RUnlock()
	return val
}

func (m *Map) Set(key, val uint64) {
	select {
	case m.items <- [2]uint64{key, val}:
	default:
	}
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

func (m *Map) process() {
	for pair := range m.items {
		m.Lock()
		m.data[pair[0]] = pair[1]
		m.Unlock()
	}
}
