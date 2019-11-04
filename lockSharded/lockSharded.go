package lockSharded

import "sync"

const (
	numShards = 256
	shardMask = numShards - 1
)

type Map struct {
	shards [numShards]lockedMap
}

func New(size uint64) *Map {
	m := &Map{}
	for i := range m.shards {
		m.shards[i].data = make(map[uint64]uint64, size/numShards)
	}
	return m
}

func (m *Map) Get(key uint64) uint64 {
	return m.shards[key&shardMask].Get(key)
}

func (m *Map) Set(key, val uint64) {
	m.shards[key&shardMask].Set(key, val)
}

func (m *Map) SetAll(pairs [][2]uint64) {
	for _, pair := range pairs {
		m.shards[pair[0]&shardMask].Set(pair[0], pair[1])
	}
}

func (m *Map) Del(key uint64) {
	m.shards[key&shardMask].Del(key)
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
