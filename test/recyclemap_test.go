/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package test

import (
    "testing"
    "github.com/xfali/goutils/container/recycleMap"
    "time"
)

func TestRecycleMap(t *testing.T) {
    dm := recycleMap.New()

    dm.Set("123", "456", time.Second)

    v1 := dm.Get("123")
    t.Logf("value is %v\n", v1)

    time.Sleep(time.Second)

    v2 := dm.Get("123")
    t.Logf("After 1 second value is %v\n", v2)

    time.Sleep(time.Millisecond)

    v3 := dm.Get("123")
    t.Logf("After 1 second 1 Millisecond value is %v\n", v3)
}
