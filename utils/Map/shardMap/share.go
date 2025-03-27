package shardMap

import (
	"fmt"
	"sync"
)

type ShardMapShared[K comparable, V any] struct {
	sync.RWMutex
	items map[K]V
}

func (m ShardMapShared[K, V]) Set(key K, value V) {
	m.Lock()
	defer m.Unlock()
	m.items[key] = value
}

type SetCb[V any] func(exist bool, valueInMap V, newValue V) V

func (m ShardMapShared[K, V]) SetCb(key K, value V, cb SetCb[V]) (res V) {
	m.Lock()
	defer m.Unlock()
	v, ok := m.items[key]
	res = cb(ok, v, value)
	m.items[key] = res
	return res
}

// /如果key不存在就设置, key存在返回false ， 不存在返回true
func (m ShardMapShared[K, V]) SetIfAbsent(key K, value V) bool {
	m.Lock()
	defer m.Unlock()
	_, ok := m.items[key]
	if !ok {
		m.items[key] = value
	}
	return !ok
}

func (m ShardMapShared[K, V]) Get(key K) (V, bool) {
	m.RLock()
	defer m.RUnlock()
	val, ok := m.items[key]
	return val, ok
}

func (m ShardMapShared[K, V]) Count() int {
	m.RLock()
	defer m.RUnlock()
	return len(m.items)
}

func (m ShardMapShared[K, V]) Has(key K) bool {
	m.RLock()
	defer m.RUnlock()
	_, ok := m.items[key]
	return ok
}

func (m ShardMapShared[K, V]) Remove(key K) {
	m.Lock()
	defer m.Unlock()
	delete(m.items, key)
}

type RemoveCb[K any, V any] func(key K, v V, exists bool) bool

func (m ShardMapShared[K, V]) RemoveCb(key K, cb RemoveCb[K, V]) bool {
	m.Lock()
	defer m.Unlock()
	v, ok := m.items[key]
	remove := cb(key, v, ok)
	if remove && ok {
		delete(shard.items, key)
	}
	return remove
}

func (m ShardMapShared[K, V]) Pop(key K) (v V, exists bool) {
	m.Lock()
	defer m.Unlock()
	v, exists = m.items[key]
	delete(shard.items, key)
	return v, exists
}

func (m ShardMapShared[K, V]) IsEmpty() bool {
	return m.Count() == 0
}
