// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package test

import (
	"github.com/xfali/goutils/container/skiplist"
	"testing"
)

func TestSkipList(t *testing.T) {
	slist := skiplist.New(skiplist.SetKeyCompareInt())
	ret := slist.Values(-1)
	for it := slist.First(); it != nil; it = it.Next() {
		t.Log("Iterator: ", it.Value())
	}

	slist.Set(-1, "s-1")
	slist.Set(1, "s1")
	slist.Set(-2, "s-2")
	slist.Set(2, "s2")
	slist.Set(3, "s3")
	if slist.Len() != 5 {
		t.Fatal("len error")
	} else {
		t.Log("len: ", slist.Len())
	}

	ret = slist.Values(-1)
	t.Log(ret)

	s2 := slist.Get(2)
	if s2.(string) != "s2" {
		t.Fatal("ret err ! ", ret[1])
	} else {
		t.Log("search: ", s2)
	}

	b := slist.Delete(2)
	if !b {
		t.Fatal("ret err ! ", ret[1])
	}

	ret = slist.Values(-1)
	t.Log("after delete s2: ", ret)

	for it := slist.First(); it != nil; it = it.Next() {
		t.Log("Iterator: ", it.Value())
	}

	slist.Set(7, "s7")
	slist.Set(10, "s10")

	for it := slist.FirstNear(7); it != nil; it = it.Next() {
		t.Log("FirstNear 7 Iterator: ", it.Value())
	}

	for it := slist.FirstNear(8); it != nil; it = it.Next() {
		t.Log("FirstNear 8 Iterator: ", it.Value())
	}

	for it := slist.FirstNear(100); it != nil; it = it.Next() {
		t.Log("FirstNear 8 Iterator: ", it.Value())
	}

	last := slist.Last()
	if last == nil {
		t.Fatal("last nil!")
	} else {
		t.Log("last value: ", last.Value())
	}

	emptyList := skiplist.New(skiplist.SetKeyCompareInt()).Last()
	if emptyList != nil {
		t.Fatal("empty list last not nil!")
	}
}
