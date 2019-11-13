package ringChanDropSharded

import "sync"

const (
	numShards  = 256
	shardMask  = numShards - 1
	stripeSize = 64
)

type Map struct {
	shards [numShards]lockedMap
	items  chan [][2]uint64
	buff   *buffer
}

func New(size uint64) *Map {
	m := &Map{
		items: make(chan [][2]uint64, size/10),
	}
	for i := range m.shards {
		m.shards[i].data = make(map[uint64]uint64, size/numShards)
	}
	go m.process()
	m.buff = newBuffer(m.Consume)
	return m
}

func (m *Map) Get(key uint64) uint64 {
	return m.shards[key&shardMask].Get(key)
}

func (m *Map) Set(key, val uint64) {
	m.buff.push([2]uint64{key, val})
}

func (m *Map) SetAll(pairs [][2]uint64) {
	for _, pair := range pairs {
		m.shards[pair[0]&shardMask].Set(pair[0], pair[1])
	}
}

func (m *Map) Del(key uint64) {
	m.shards[key&shardMask].Del(key)
}

func (m *Map) Consume(pairs [][2]uint64) {
	select {
	case m.items <- pairs:
	default:
	}
}

func (m *Map) process() {
	for items := range m.items {
		for _, item := range items {
			m.shards[item[0]&shardMask].Set(item[0], item[1])
		}
	}
}

////////////////////////////////////////////////////////////////////////////////

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

////////////////////////////////////////////////////////////////////////////////

type buffer struct {
	pool *sync.Pool
}

func newBuffer(consume func([][2]uint64)) *buffer {
	return &buffer{
		pool: &sync.Pool{
			New: func() interface{} { return newBufferStripe(consume) },
		},
	}
}

func (b *buffer) push(pair [2]uint64) {
	stripe := b.pool.Get().(*bufferStripe)
	stripe.push(pair)
	b.pool.Put(stripe)
}

type bufferStripe struct {
	consume func([][2]uint64)
	data    [][2]uint64
}

func newBufferStripe(consume func([][2]uint64)) *bufferStripe {
	return &bufferStripe{
		consume: consume,
		data:    make([][2]uint64, 0, stripeSize),
	}
}

func (s *bufferStripe) push(pair [2]uint64) {
	s.data = append(s.data, pair)
	if len(s.data) >= cap(s.data) {
		s.consume(s.data)
		s.data = s.data[:0]
	}
}
