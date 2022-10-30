package stc

import (
	"time"

	"go.oneofone.dev/genh"
)

type entry[V any] struct {
	val V
	t   *time.Timer
	exp int64
}

func (e *entry[V]) cancel() {
	if e == nil || e.t == nil {
		return
	}
	if !e.t.Stop() {
		select {
		case <-e.t.C:
		default:
		}
	}
}

type SimpleTimedCache[K comparable, V any] struct {
	m genh.LMap[K, *entry[V]]

	OnSet func(key K, val V)
}

func (stc *SimpleTimedCache[K, V]) Set(key K, val V, expiry time.Duration) {
	expireAt := time.Now().Add(expiry)
	stc.m.DeleteGet(key).cancel()

	e := &entry[V]{val: val, exp: expireAt.Unix()}
	if expiry > 0 {
		e.t = time.AfterFunc(expiry, func() {
			stc.m.Update(func(m map[K]*entry[V]) {
				if e := m[key]; e != nil && time.Until(expireAt) < 1 {
					delete(m, key)
				}
			})
		})
	}
	stc.m.Set(key, e)
	if stc.OnSet != nil {
		stc.OnSet(key, val)
	}
}

func (stc *SimpleTimedCache[K, V]) SetWithUpdate(key K, valFn func() V, every time.Duration) {
	val := valFn()
	stc.m.DeleteGet(key).cancel()

	e := &entry[V]{val: val, exp: -1}
	if every > 1 {
		e.t = time.AfterFunc(every, func() {
			if e := stc.m.Get(key); e == nil || e.exp != -1 {
				return
			}
			stc.SetWithUpdate(key, valFn, every)
		})
	}
	stc.m.Set(key, e)

	if stc.OnSet != nil {
		stc.OnSet(key, val)
	}
}

func (stc *SimpleTimedCache[K, V]) Delete(key K) {
	stc.m.DeleteGet(key).cancel()
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
