package stc

import (
	"fmt"
	"testing"
	"time"
)

func TestLeak(t *testing.T) {
	var stc SimpleTimedCache[string, string]
	stc.OnSet = func(key, val string) {}
	i := 0
	stc.SetFnWithExpire("key", func() string {
		i++
		return fmt.Sprintf("val:%d", i)
	}, time.Millisecond*100, time.Millisecond*400)
	time.Sleep(time.Millisecond * 100)
	if v := stc.Get("key"); v != "val:1" {
		t.Fatal("expected val:1, got", v)
	}
	time.Sleep(time.Millisecond * 350)
	if v := stc.Get("key"); v != "val:5" {
		t.Fatal("expected val:5, got", v)
	}
	time.Sleep(time.Millisecond * 500)
	if v := stc.Get("key"); v != "" {
		t.Fatal("expected empty, got", v)
	}

	stc.Set("key2", "val1", time.Millisecond*100)
	time.Sleep(time.Millisecond * 350)
	if v := stc.Get("key2"); v != "" {
		t.Fatal("expected empty, got", v)
	}
	stc.Set("key2", "val2", time.Minute)
	time.Sleep(time.Millisecond * 350)
	stc.Set("key2", "val3", time.Minute)

	if v, ok := stc.GetOk("key2"); v != "val3" || !ok {
		t.Fatalf("unexpected val3: %v", v)
	}
}
