package stc

import (
	"time"

	"go.oneofone.dev/genh"
)

type entry[V any] struct {
	val V
	t   *time.Timer
	exp bool
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
	m       genh.LMap[K, *entry[V]]
	lastUse genh.LMap[K, int64]

	OnSet func(key K, val V)
}

func (stc *SimpleTimedCache[K, V]) Set(key K, val V, expiry time.Duration) {
	expireAt := time.Now().Add(expiry)
	stc.m.DeleteGet(key).cancel()

	e := &entry[V]{val: val, exp: expiry > 0}
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

func (stc *SimpleTimedCache[K, V]) SetFn(key K, valFn func() V, every time.Duration) {
	stc.SetFnWithExpire(key, valFn, every, 0)
}

func (stc *SimpleTimedCache[K, V]) SetFnWithExpire(key K, valFn func() V, every, expireIfNotUsedAfter time.Duration) {
	val := valFn()
	stc.m.DeleteGet(key).cancel()

	e := &entry[V]{val: val}
	if every > 1 {
		e.t = time.AfterFunc(every, func() {
			if expireIfNotUsedAfter > 0 {
				lu := stc.lastUse.Get(key)
				if time.Since(time.Unix(0, lu)) > expireIfNotUsedAfter {
					stc.Delete(key)
					return
				}
			}
			stc.SetFnWithExpire(key, valFn, every, expireIfNotUsedAfter)
		})
	}
	stc.m.Set(key, e)

	if stc.OnSet != nil {
		stc.OnSet(key, val)
	}
}

func (stc *SimpleTimedCache[K, V]) Delete(key K) {
	stc.lastUse.Delete(key)
	stc.m.DeleteGet(key).cancel()
}

func (stc *SimpleTimedCache[K, V]) Get(key K) (_ V) {
	if e := stc.m.Get(key); e != nil {
		stc.lastUse.Set(key, time.Now().UnixNano())
		return e.val
	}
	return
}

func (stc *SimpleTimedCache[K, V]) GetOk(key K) (_ V, ok bool) {
	if e := stc.m.Get(key); e != nil {
		stc.lastUse.Set(key, time.Now().Unix())
		return e.val, true
	}
	return
}
