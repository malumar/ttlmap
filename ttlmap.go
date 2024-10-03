package ttlmap

import (
	"sync"
	"time"
)

type item[T any] struct {
	Value      T
	lastAccess int64
}

type TTLMap[K comparable, T any] struct {
	m      map[K]*item[T]
	l      sync.Mutex
	closer bool
}

func New[K comparable, T any](ln int, maxTTL int, closer func(val *item[T])) (m *TTLMap[K, T]) {
	m = &TTLMap[K, T]{m: make(map[K]*item[T], ln)}
	if closer != nil {
		go func() {
			for now := range time.Tick(time.Second) {
				m.l.Lock()
				for k, v := range m.m {
					if now.Unix()-v.lastAccess > int64(maxTTL) {
						closer(v)
						delete(m.m, k)
					}
				}
				m.l.Unlock()
			}
		}()

	} else {
		go func() {
			for now := range time.Tick(time.Second) {
				m.l.Lock()
				for k, v := range m.m {
					if now.Unix()-v.lastAccess > int64(maxTTL) {
						delete(m.m, k)
					}
				}
				m.l.Unlock()
			}
		}()

	}
	return
}

func (m *TTLMap[K, T]) Len() int {
	return len(m.m)
}

func (m *TTLMap[K, T]) Put(k K, v T) {
	m.l.Lock()
	it, ok := m.m[k]
	if !ok {
		it = &item[T]{Value: v}
		m.m[k] = it
	}
	it.lastAccess = time.Now().Unix()
	m.l.Unlock()
}

func (m *TTLMap[K, T]) Get(k K) (v T, found bool) {
	m.l.Lock()
	if it, ok := m.m[k]; ok {
		v = it.Value
		found = true
		it.lastAccess = time.Now().Unix()
	}
	m.l.Unlock()
	return

}
