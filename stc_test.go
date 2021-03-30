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

	stc.Set("key2", "val1", time.Millisecond*10)
	stc.Set("key2", "val2", time.Minute)
	time.Sleep(time.Millisecond * 100)
	if stc.Get("key2") != "val2" {
		t.Fatalf("unexpected val: %v", stc.Get("key2"))
	}
}
