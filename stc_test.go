package stc

import (
	"testing"
	"time"
)

func TestSTC(t *testing.T) {
	var stc SimpleTimedCache[string, string]
	stc.OnSet = func(key, val string) {}
	stc.Set("a", "b", time.Second)
	stc.Set("b", "b", time.Second*5)
	if stc.Get("a") != "b" {
		t.Fatal("get failed")
	}
	time.Sleep(time.Second * 2)
	if stc.Get("a") != "" {
		t.Fatal("get failed")
	}
	if stc.Get("b") != "b" {
		t.Fatal("get failed")
	}
}
