package stc

import (
	"sync"
	"time"
)

type SimpleTimedCache struct {
	m sync.Map
}

func (stc *SimpleTimedCache) Set(key, val interface{}, expiry time.Duration) {
	stc.m.Store(key, val)
	if expiry > 0 {
		time.AfterFunc(expiry, func() {
			if cv, _ := stc.m.Load(key); cv == val {
				stc.m.Delete(key)
			}
		})
	}
}

func (stc *SimpleTimedCache) SetWithUpdate(key interface{}, valFn func() interface{}, every time.Duration) {
	val := valFn()
	stc.m.Store(key, val)
	if every == 0 {
		return
	}

	var fn func()
	fn = func() {
		if !stc.Exists(key) {
			return
		}
		stc.m.Store(key, valFn())
		time.AfterFunc(every, fn)
	}
	time.AfterFunc(every, fn)
}

func (stc *SimpleTimedCache) Delete(key interface{}) {
	stc.m.Delete(key)
}

func (stc *SimpleTimedCache) Exists(key interface{}) bool {
	_, ok := stc.m.Load(key)
	return ok
}

func (stc *SimpleTimedCache) Get(key interface{}) interface{} {
	if v, ok := stc.GetOk(key); ok {
		return v
	}
	return nil
}

func (stc *SimpleTimedCache) GetOk(key interface{}) (interface{}, bool) {
	return stc.m.Load(key)
}
