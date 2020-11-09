// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package test

import (
	"github.com/xfali/goutils/container/lru"
	"testing"
)

func TestSimpleLRU(t *testing.T) {
	lru := lru.NewLruCache(3)
	testLru(t, lru)
}

func testLru(t *testing.T, m lru.LRU) {
	if _, ok := m.Get(1); ok {
		t.Fatal("key 1 have no value ")
	}

	m.Put(1, "a")
	if _, ok := m.Get(2); ok {
		t.Fatal("key 2 have no value ")
	}
	if v, ok := m.Get(1); !ok || v.(string) != "a" {
		t.Fatal("not exits, v: ", v)
	}

	m.Delete(2)

	m.Put(2, "b")
	if v, ok := m.Get(2); !ok || v.(string) != "b" {
		t.Fatal("not exits, v: ", v)
	}

	if m.Size() != 2 {
		t.Fatal("must 2 but: ", m.Size())
	}

	m.Delete(2)

	if v, ok := m.Get(2); ok {
		t.Fatal("2 must be not exists, v: ", v)
	}

	if m.Size() != 1 {
		t.Fatal("must 2 but: ", m.Size())
	}

	//test elimination
	// 1 exists
	m.Put(2, "b")
	// 1 exists
	m.Put(3, "c")
	// 1 eliminated
	m.Put(4, "d")
	if v, ok := m.Get(1); ok {
		t.Fatal("1 must not exits, v: ", v)
	}
}
