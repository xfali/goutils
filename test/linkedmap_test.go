// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package test

import (
	"github.com/xfali/goutils/container/xmap"
	"testing"
)

func TestMap(t *testing.T) {
	m := xmap.NewSimpleMap()
	testMap(t, m)
}

func TestSimpleLinkedMap(t *testing.T) {
	m := xmap.NewSimpleLinkedMap()
	testMap(t, m)

	m.Put(3, "c")
	m.Foreach(func(key interface{}, value interface{}) bool {
		if key.(int) != 1 || value.(string) != "a" {
			t.Fatal("key ", key, " value ", value, "not match")
		}
		return false
	})
}

func TestLinkedMap(t *testing.T) {
	m := xmap.NewLinkedMap()
	testMap(t, m)

	m.Put(3, "c")
	m.Foreach(func(key interface{}, value interface{}) bool {
		if key.(int) != 1 || value.(string) != "a" {
			t.Fatal("key ", key, " value ", value, "not match")
		}
		return false
	})
}

func testMap(t *testing.T, m xmap.Map) {
	if _, ok := m.Get(1) ; ok {
		t.Fatal("key 2 have no value ")
	}

	m.Put(1, "a")
	if _, ok := m.Get(2) ; ok {
		t.Fatal("key 2 have no value ")
	}
	if v, ok := m.Get(1) ; !ok || v.(string) != "a" {
		t.Fatal("not exits, v: ", v)
	}

	if !m.Find(1) {
		t.Fatal("cannot find 1?")
	}

	v, load := m.GetOrPut(1, "x")
	if !load {
		t.Fatal("must loaded")
	}
	if v.(string) != "a" {
		t.Fatal("must be a but: ", v)
	}

	v, load = m.GetOrPut(2, "b")
	if load {
		t.Fatal("must not loaded")
	}
	if v.(string) != "b" {
		t.Fatal("must be b but: ", v)
	}

	m.Foreach(func(key interface{}, value interface{}) bool {
		t.Log("key ", key, " value ", value)
		return true
	})

	m.Delete(2)
	if m.Find(2) {
		t.Fatal("cannot find 2")
	}

	if m.Size() != 1 {
		t.Fatal("must 1 but: ", m.Size())
	}
}
