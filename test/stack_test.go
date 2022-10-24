// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package test

import (
	"github.com/xfali/goutils/v2/container/stack"
	"testing"
)

func TestStack(t *testing.T) {
	s := stack.New()

	s.Push("test1")
	s.Push("test2")

	if s.Len() != 2 {
		t.Fatal("s len error")
	}

	t.Log("len: ", s.Len())

	s.Foreach(func(i interface{}) bool {
		t.Log("foreach: ", i)
		return false
	})

	v := s.Pop().(string)
	if v != "test2" {
		t.Fatal("value error: ", v)
	}

	t.Log("pop: ", v)

	if s.Len() != 1 {
		t.Fatal("s len error")
	}

	t.Log("len: ", s.Len())
}
