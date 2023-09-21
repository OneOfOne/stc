package stc

import (
	"testing"
	"time"
)

func TestSetGet(t *testing.T) {
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

func TestMustGet(t *testing.T) {
	var stc SimpleTimedCache[string, string]
	_ = stc.MustGet("a", func() string { return "a" }, time.Second)
	if v := stc.MustGet("a", func() string { return "b" }, time.Second); v != "a" {
		t.Fatal("must get failed", v)
	}
	time.Sleep(time.Second * 2)
	if v := stc.MustGet("a", func() string { return "c" }, time.Second); v != "c" {
		t.Fatal("must get failed", v)
	}
}
