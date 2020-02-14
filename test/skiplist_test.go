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

	slist.Insert(1, "s1")
	slist.Insert(2, "s2")
	slist.Insert(3, "s3")
	if slist.Len() != 3 {
		t.Fatal("len error")
	} else {
		t.Log("len: ", slist.Len())
	}

	ret = slist.Values(-1)
	if ret[0].(string) != "s1" {
		t.Fatal("ret err ! ", ret[0])
	}
	if ret[1].(string) != "s2" {
		t.Fatal("ret err ! ", ret[1])
	}
	if ret[2].(string) != "s3" {
		t.Fatal("ret err ! ", ret[2])
	}

	s2 := slist.Search(2)
	if s2.(string) != "s2" {
		t.Fatal("ret err ! ", ret[1])
	}

	b := slist.Delete(2)
	if !b {
		t.Fatal("ret err ! ", ret[1])
	}

	ret = slist.Values(-1)
	for _, v := range ret {
		if v.(string) == s2 {
			t.Fatal()
		}
		t.Log("after delete")
		t.Log(v)
	}
}
