// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package test

import (
	"github.com/xfali/goutils/v2/container/linkedSet"
	"github.com/xfali/goutils/v2/container/xlist"
	"testing"
)

func TestSimpleList(t *testing.T) {
	l := xlist.NewSimpleList()
	testList(t, l)
}

func TestLinkedSet(t *testing.T) {
	l := linkedSet.New()
	testList(t, l)
}

func testList(t *testing.T, m xlist.List) {
	m.PushBack("a")
	v := m.Front()
	if v.(string) != "a" {
		t.Fatal("not match, expect a but get: ", v)
	}

	m.PushFront("b")
	v = m.Front()
	if v.(string) != "b" {
		t.Fatal("not match, expect b but get: ", v)
	}

	v = m.Back()
	if v.(string) != "a" {
		t.Fatal("not match, expect a but get: ", v)
	}

	if m.Len() != 2 {
		t.Fatal("not match, expect 2 but get: ", m.Len())
	}

	m.Foreach(func(i interface{}) bool {
		t.Log(i)
		return true
	})

	v = m.PopBack()
	if v.(string) != "a" {
		t.Fatal("not match, expect a but get: ", v)
	}

	if m.Len() != 1 {
		t.Fatal("not match, expect 1 but get: ", m.Len())
	}

	b := m.Find("c")
	if b {
		t.Fatal("not match, expect false but get: ", b)
	}

	m.PushBack("c")

	b = m.Find("c")
	if !b {
		t.Fatal("not match, expect true but get: ", b)
	}

	v = m.PopFront()
	if v.(string) != "b" {
		t.Fatal("not match, expect b but get: ", v)
	}

	m.Remove("c")
	if m.Len() != 0 {
		t.Fatal("not match, expect 0 but get: ", m.Len())
	}

	m.Foreach(func(i interface{}) bool {
		t.Log(i)
		return true
	})
}
