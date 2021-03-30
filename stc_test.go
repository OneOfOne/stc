package stc

import (
	"runtime"
	"testing"
	"time"
)

func TestLeak(t *testing.T) {
	var stc SimpleTimedCache
	stc.SetWithUpdate("key", func() interface{} {
		t.Log(runtime.NumGoroutine(), stc.Get("key"))
		return "val"
	}, time.Millisecond*100)
	t.Log(runtime.NumGoroutine(), stc.Get("key"))
	time.Sleep(time.Millisecond * 350)
	stc.Delete("key")
	t.Log(runtime.NumGoroutine(), stc.Get("key"))
	time.Sleep(time.Millisecond * 150)
	t.Log(runtime.NumGoroutine(), stc.Get("key"))
}
