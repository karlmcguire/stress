package syncChanDropSharded

import "sync"

const (
	numShards = 256
	shardMask = numShards - 1
)

type Map struct {
	shards [numShards]*sync.Map
	items  chan [3]uint64
}

func New(size uint64) *Map {
	m := &Map{
		items: make(chan [3]uint64, size/10),
	}
	for i := range m.shards {
		m.shards[i] = &sync.Map{}
	}
	go m.process()
	return m
}

func (m *Map) Get(key uint64) uint64 {
	val, _ := m.shards[key&shardMask].Load(key)
	if val == nil {
		return 0
	}
	return val.(uint64)
}

func (m *Map) Set(key, val uint64) {
	select {
	case m.items <- [3]uint64{key & shardMask, key, val}:
	default:
	}
}

func (m *Map) SetAll(pairs [][2]uint64) {
	for _, pair := range pairs {
		m.shards[pair[0]&shardMask].Store(pair[0], pair[1])
	}
}

func (m *Map) Del(key uint64) {
	m.shards[key&shardMask].Delete(key)
}

func (m *Map) process() {
	for item := range m.items {
		m.shards[item[0]].Store(item[1], item[2])
	}
}
