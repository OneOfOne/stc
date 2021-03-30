package stc

import (
	"sync"
	"time"
)

type entry struct {
	val      interface{}
	expireAt int64
}

type SimpleTimedCache struct {
	m   map[string]*entry
	mux sync.RWMutex
}

func (stc *SimpleTimedCache) init() *SimpleTimedCache {
	if stc.m == nil {
		stc.m = make(map[string]*entry)
	}
	return stc
}

func (stc *SimpleTimedCache) Set(key string, val interface{}, expiry time.Duration) {
	expireAt := time.Now().Add(expiry).Unix()
	stc.mux.Lock()
	stc.init().m[key] = &entry{val, expireAt}
	stc.mux.Unlock()
	if expiry > 0 {
		time.AfterFunc(expiry, func() {
			stc.mux.Lock()
			if e := stc.m[key]; e != nil && e.expireAt == expireAt {
				delete(stc.m, key)
			}
			stc.mux.Unlock()
		})
	}
}

func (stc *SimpleTimedCache) SetWithUpdate(key string, valFn func() interface{}, every time.Duration) {
	val := valFn()
	stc.mux.Lock()
	stc.init().m[key] = &entry{val, -1}
	stc.mux.Unlock()
	if every == 0 {
		return
	}

	var fn func()
	fn = func() {
		e := stc.get(key)
		if e == nil || e.expireAt != -1 {
			return
		}
		val := valFn()
		stc.mux.Lock()
		e = stc.m[key]
		if e != nil && e.expireAt == -1 {
			e.val = val
		}
		stc.mux.Unlock()
		if e != nil && e.expireAt == -1 {
			time.AfterFunc(every, fn)
		}
	}
	time.AfterFunc(every, fn)
}

func (stc *SimpleTimedCache) Delete(key string) {
	stc.mux.Lock()
	delete(stc.m, key)
	stc.mux.Unlock()
}

func (stc *SimpleTimedCache) Exists(key string) bool {
	stc.mux.RLock()
	ok := stc.m[key] != nil
	stc.mux.RUnlock()
	return ok
}

func (stc *SimpleTimedCache) Get(key string) interface{} {
	if e := stc.get(key); e != nil {
		return e.val
	}
	return nil
}

func (stc *SimpleTimedCache) GetOk(key string) (interface{}, bool) {
	if e := stc.get(key); e != nil {
		return e.val, true
	}
	return nil, false
}

func (stc *SimpleTimedCache) get(key string) *entry {
	stc.mux.RLock()
	e := stc.m[key]
	stc.mux.RUnlock()
	return e
}
