package sync

import "sync"

type Map struct {
	data *sync.Map
}

func New(size uint64) *Map {
	return &Map{
		data: &sync.Map{},
	}
}

func (m *Map) Get(key uint64) uint64 {
	val, _ := m.data.Load(key)
	if val == nil {
		return 0
	}
	return val.(uint64)
}

func (m *Map) Set(key, val uint64) {
	m.data.Store(key, val)
}

func (m *Map) SetAll(pairs [][2]uint64) {
	for _, pair := range pairs {
		m.data.Store(pair[0], pair[1])
	}
}

func (m *Map) Del(key uint64) {
	m.data.Delete(key)
}
