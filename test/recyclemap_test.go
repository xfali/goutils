/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description:
 */

package test

import (
	"github.com/xfali/goutils/container/recycleMap"
	"testing"
	"time"
)

func TestRecycleMap(t *testing.T) {
	dm := recycleMap.New()
	test(dm, t)
}

func TestRecycleMapWithNotifier(t *testing.T) {
	dm := recycleMap.New(recycleMap.OptSetDeleteNotifier(func(key, value interface{}) {
		t.Logf(">>>>>>>>>>>>>>>> key deleted: %v , value: %v\n", key, value)
	}))
	test(dm, t)
}

func test(dm recycleMap.RecycleMap, t *testing.T) {
	defer dm.Close()

	dm.Set("123", "abc", 50*time.Millisecond)
	dm.Set("456", "efg", -1)
	dm.Set("789", "hij", -1)

	v1 := dm.Get("123")
	if v1.(string) != "abc" {
		t.Fatal("must be abc")
	}
	t.Logf("value is %v, ttl: %d\n", v1, dm.TTL("123")/time.Millisecond)
	size := dm.Size()
	if size != 3 {
		t.Fatal("expect 3 but get ", size)
	}

	time.Sleep(50 * time.Millisecond)

	v2 := dm.Get("123")
	if v2 != nil {
		t.Fatalf("v2 must be nil, %d\n", dm.TTL("123")/time.Millisecond)
	}
	t.Logf("After 1 second value is %v\n", v2)
	size = dm.Size()
	if size != 2 {
		t.Fatal("expect 3 but get ", size)
	}

	time.Sleep(time.Millisecond)

	v3 := dm.Get("123")
	if v3 != nil {
		t.Fatalf("v2 must be nil, %d\n", dm.TTL("123")/time.Millisecond)
	}
	t.Logf("After 1 second 1 Millisecond value is %v\n", v3)

	v4 := dm.Get("456")
	if v4.(string) != "efg" {
		t.Fatal("must be efg")
	}
	t.Log(v4, dm.TTL("456"))

	dm.SetExpire("456", 50*time.Millisecond)
	t.Log(dm.TTL("456"))
	time.Sleep(50 * time.Millisecond)
	v5 := dm.Get("456")
	if v5 != nil {
		t.Fatal("must be nil")
	}
	size = dm.Size()
	if size != 1 {
		t.Fatal("expect 1 but get ", size)
	}

	v6 := dm.Get("789")
	if v6.(string) != "hij" {
		t.Fatal("must be hij")
	} else {
		t.Log(v6)
	}

	dm.Delete("789")
	if dm.Get("789") != nil {
		t.Fatal("must be nil")
	}

	time.Sleep(50 * time.Millisecond)
	size = dm.Size()
	if size != 0 {
		t.Fatal("expect 3 but get ", size)
	}
}
