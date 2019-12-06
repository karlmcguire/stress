package shardedPolicy

import (
	"sync"
)

const (
	numShards = 256
	shardMask = numShards - 1
)

type Map struct {
	shards [numShards]shard
}

func New(size uint64) *Map {
	m := &Map{}
	for i := range m.shards {
		m.shards[i].data = make(map[uint64]item, size/numShards)
	}
	return m
}

func (m *Map) Get(key uint64) uint64 {
	return m.shards[key&shardMask].Get(key)
}

func (m *Map) Set(key, val uint64, cost, ttl int64) {
	m.shards[key&shardMask].Set(key, val, cost, ttl)
}

func (m *Map) Del(key uint64) {
	m.shards[key&shardMask].Del(key)
}

////////////////////////////////////////////////////////////////////////////////

type shard struct {
	sync.RWMutex
	data map[uint64]item
}

type item struct {
	cost int64
	ttl  int64
	val  uint64
}

func (s *shard) Get(key uint64) uint64 {
	s.RLock()
	data := s.data[key]
	s.RUnlock()
	return data.val
}

func (s *shard) Set(key, val uint64, cost, ttl int64) {
	s.Lock()
	s.data[key] = item{val: val, cost: cost, ttl: ttl}
	s.Unlock()
}

func (s *shard) Del(key uint64) {
	s.Lock()
	delete(s.data, key)
	s.Unlock()
}
