/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description:
 */

package test

import (
	"github.com/xfali/goutils/v2/container/recycleMap"
	"github.com/xfali/goutils/v2/pattern"
	"sync"
	"testing"
	"time"
)

func TestRecycleMap(t *testing.T) {
	dm := recycleMap.New[string, string]()
	t.Run("test", func(t *testing.T) {
		test(dm, t)
	})
}

func TestRecycleMapKeys(t *testing.T) {
	t.Run("keys", func(t *testing.T) {
		dm := recycleMap.New[string, string]()
		dm.Set("aaabbbccc", "111222333", -1)
		dm.Set("bbbaaaccc", "222111333", -1)
		dm.Set("cccaaabbb", "333111222", -1)
		keys := dm.Keys("")
		for _, k := range keys {
			t.Log("keys: ", k)
		}
	})

	t.Run("keys regexp", func(t *testing.T) {
		dm := recycleMap.New[string, string](
			recycleMap.OptSetMatcher[string, string](recycleMap.RegexpMatcher(func(v string) string {
				return v
			})))
		dm.Set("aaabbbccc", "111222333", -1)
		dm.Set("bbbaaaccc", "222111333", -1)
		dm.Set("cccaaabbb", "333111222", -1)
		keys := dm.Keys(pattern.StartWith("aaa"))
		for _, k := range keys {
			t.Log("keys: ", k)
		}
		if len(keys) != 1 {
			t.Fatal("expect 1 but get ", len(keys))
		}
		if keys[0] != "aaabbbccc" {
			t.Fatal("expect aaabbbccc but get ", keys[0])
		}
	})
}

func TestRecycleMapWithNotifier(t *testing.T) {
	t.Run("1 map", func(t *testing.T) {
		dm := recycleMap.New(recycleMap.OptSetDeleteNotifier(func(key, value string) {
			t.Logf(">>>>>>>>>>>>>>>> key deleted: %v , value: %v\n", key, value)
		}))
		test(dm, t)
	})

	t.Run("multi map", func(t *testing.T) {
		total := 10
		wait := sync.WaitGroup{}

		dms := make([]recycleMap.RecycleMap[string, string], total)
		for i := 0; i < total; i++ {
			dm := recycleMap.New(recycleMap.OptSetDeleteNotifier(func(key, value string) {
				t.Logf(">>>>>>>>>>>>>>>> key deleted: %v , value: %v\n", key, value)
			}))
			dms[i] = dm
			wait.Add(1)
		}
		for i := range dms {
			dm := dms[i]
			go func() {
				defer wait.Done()
				test(dm, t)
			}()
		}
		wait.Wait()
	})
}

func test(dm recycleMap.RecycleMap[string, string], t *testing.T) {
	defer dm.Close()

	dm.Set("123", "abc", 50*time.Millisecond)
	dm.Set("456", "efg", -1)
	dm.Set("789", "hij", -1)

	keys := dm.Keys("")
	for _, k := range keys {
		t.Log("keys: ", k)
	}

	v1 := dm.Get("123")
	if v1 != "abc" {
		t.Fatal("must be abc")
	}
	t.Logf("value is %v, ttl: %d\n", v1, dm.TTL("123")/time.Millisecond)
	size := dm.Size()
	if size != 3 {
		t.Fatal("expect 3 but get ", size)
	}

	time.Sleep(2 * time.Second)
	if dm.Size() != 2 {
		t.Fatal("expect 2 but get ", dm.Size())
	}
	v2 := dm.Get("123")
	if v2 != "" {
		t.Fatalf("v2 must be nil, %d\n", dm.TTL("123")/time.Millisecond)
	}
	t.Logf("After 2 second value is %v\n", v2)
	size = dm.Size()
	if size != 2 {
		t.Fatal("expect 3 but get ", size)
	}

	time.Sleep(time.Millisecond)

	v3 := dm.Get("123")
	if v3 != "" {
		t.Fatalf("v2 must be nil, %d\n", dm.TTL("123")/time.Millisecond)
	}
	t.Logf("After 1 second 1 Millisecond value is %v\n", v3)

	v4 := dm.Get("456")
	if v4 != "efg" {
		t.Fatal("must be efg")
	}
	t.Log(v4, dm.TTL("456"))

	dm.SetExpire("456", 50*time.Millisecond)
	t.Log(dm.TTL("456"))
	time.Sleep(50 * time.Millisecond)
	v5 := dm.Get("456")
	if v5 != "" {
		t.Fatal("must be nil")
	}
	size = dm.Size()
	if size != 1 {
		t.Fatal("expect 1 but get ", size)
	}

	v6 := dm.Get("789")
	if v6 != "hij" {
		t.Fatal("must be hij")
	} else {
		t.Log(v6)
	}

	dm.Delete("789")
	if dm.Get("789") != "" {
		t.Fatal("must be nil")
	}

	time.Sleep(50 * time.Millisecond)
	size = dm.Size()
	if size != 0 {
		t.Fatal("expect 3 but get ", size)
	}
}
