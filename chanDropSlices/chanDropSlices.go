package chanDropSlices

import "sync"

const (
	numShards = 256
	shardMask = numShards - 1
)

type Map struct {
	shards [numShards]lockedMap
	items  chan [][2]uint64
}

func New(size uint64) *Map {
	m := &Map{
		items: make(chan [][2]uint64, size/10),
	}
	for i := range m.shards {
		m.shards[i].data = make(map[uint64]uint64, size/numShards)
	}
	go m.process()
	return m
}

func (m *Map) Get(key uint64) uint64 {
	return m.shards[key&shardMask].Get(key)
}

func (m *Map) Set(key, val uint64) {
	select {
	case m.items <- [][2]uint64{[2]uint64{key, val}}:
	default:
	}
}

func (m *Map) SetAll(items [][2]uint64) {
	for _, item := range items {
		m.shards[item[0]&shardMask].Set(item[0], item[1])
	}
}

func (m *Map) Del(key uint64) {
	m.shards[key&shardMask].Del(key)
}

func (m *Map) process() {
	for items := range m.items {
		for _, item := range items {
			m.shards[item[0]&shardMask].Set(item[0], item[1])
		}
	}
}

type lockedMap struct {
	sync.RWMutex
	data map[uint64]uint64
}

func (m *lockedMap) Get(key uint64) uint64 {
	m.RLock()
	val := m.data[key]
	m.RUnlock()
	return val
}

func (m *lockedMap) Set(key, val uint64) {
	m.Lock()
	m.data[key] = val
	m.Unlock()
}

func (m *lockedMap) Del(key uint64) {
	m.Lock()
	delete(m.data, key)
	m.Unlock()
}
