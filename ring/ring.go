package ring

import "sync"

const stripeSize = 64

type Map struct {
	sync.RWMutex
	data map[uint64]uint64
	buff *buffer
}

func New(size uint64) *Map {
	m := &Map{
		data: make(map[uint64]uint64, size),
	}
	m.buff = newBuffer(m.Consume)
	return m
}

func (m *Map) Get(key uint64) uint64 {
	m.RLock()
	val := m.data[key]
	m.RUnlock()
	return val
}

func (m *Map) Set(key, val uint64) {
	m.buff.push([2]uint64{key, val})
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

func (m *Map) Consume(pairs [][2]uint64) {
	m.Lock()
	for _, pair := range pairs {
		m.data[pair[0]] = pair[1]
	}
	m.Unlock()
}

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
