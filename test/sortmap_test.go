// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package test

import (
	"github.com/xfali/goutils/container/sortmap"
	"testing"
	"time"
)

func TestSortMap(t *testing.T) {
	sm := sortmap.New()
	err := sm.Add("int", 1, "time", time.Now(), "nil")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("keys", sm.Keys())
	t.Log("all", sm.GetAll())
	it := sm.Iterator()
	for it.HasNext() {
		x, v := it.Next()
		t.Log("iterator", x, v)
	}

	sm.Add("float", 1.1, "string", "test")
	t.Log("after add twice")
	it = sm.Iterator()
	for it.HasNext() {
		x, v := it.Next()
		t.Log("iterator", x, v)
	}

	sm.Remove("float")
	t.Log("after remove float")
	it = sm.Iterator()
	for it.HasNext() {
		x, v := it.Next()
		t.Log("iterator", x, v)
	}
}
