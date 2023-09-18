package stc

import (
	"sync"
	"time"

	"go.oneofone.dev/genh"
)

type entry[V any] struct {
	val     V
	created int64
	ttl     time.Duration
}

type SimpleTimedCache[K comparable, V any] struct {
	m     genh.LMap[K, *entry[V]]
	OnSet func(key K, val V)
	init  sync.Once
}

func (stc *SimpleTimedCache[K, V]) Set(key K, val V, ttl time.Duration) {
	stc.init.Do(func() {
		go stc.cleanup()
	})
	e := &entry[V]{val: val, created: time.Now().UnixNano(), ttl: ttl}
	stc.m.Set(key, e)
	if stc.OnSet != nil {
		stc.OnSet(key, val)
	}
}

func (stc *SimpleTimedCache[K, V]) Delete(key K) {
	stc.m.Delete(key)
}

func (stc *SimpleTimedCache[K, V]) MustGet(key K, fn func() V, ttl time.Duration) (_ V) {
	stc.init.Do(func() {
		go stc.cleanup()
	})
	return stc.m.MustGet(key, func() *entry[V] {
		return &entry[V]{val: fn(), created: time.Now().Unix(), ttl: ttl}
	}).val
}

func (stc *SimpleTimedCache[K, V]) Get(key K) (_ V) {
	if e := stc.m.Get(key); e != nil {
		return e.val
	}
	return
}

func (stc *SimpleTimedCache[K, V]) GetOk(key K) (_ V, ok bool) {
	if e := stc.m.Get(key); e != nil {
		return e.val, true
	}
	return
}

func (stc *SimpleTimedCache[K, V]) cleanup() {
	for {
		time.Sleep(time.Second)
		var keys []K
		now := time.Now().UnixNano()
		stc.m.ForEach(func(key K, e *entry[V]) bool {
			if e.ttl > 0 && now-e.created > int64(e.ttl) {
				keys = append(keys, key)
			}
			return true
		})

		for _, k := range keys {
			stc.m.Delete(k)
		}
	}
}
