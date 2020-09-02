/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description:
 */

package test

import (
	"container/list"
	"fmt"
	"github.com/xfali/goutils/idUtil"
	"testing"
	"time"
)

func TestRandomId(t *testing.T) {
	t.Logf("%s\n", idUtil.RandomId(10))
	t.Logf("%s\n", idUtil.RandomId(32))
	t.Logf("%s\n", idUtil.RandomId(64))
}

func TestSnowFlakeId(t *testing.T) {
	sf := idUtil.NewSnowFlake()
	id, err := sf.NextId()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("idUtil is %d, timestamp %d limitstr %s str %s\n", id.Int64(), id.Timestamp(), id.LimitString(30), id.String())
	fmt.Printf("time is %v\n", id.Time())

	str := idUtil.Compress2StringUL2(id.Int64(), 20)
	fmt.Printf("func compress str is %s\n", str)
	id2 := idUtil.Uncompress2LongUL(str)
	fmt.Printf("func uncompress %d\n", id2)

	sid := id.Compress()
	fmt.Printf("compress id %s\n", sid)
	fmt.Printf("uncompress id %d\n", sid.UnCompress())
}

func TestSnowFlakeId2(t *testing.T) {
	sf := idUtil.NewSnowFlake()
	i := 0
	now := time.Now()
	l := list.New()
	for {
		id, _ := sf.NextId()
		l.PushBack(id)
		if i == 10000 {
			break
		}
		i++
	}
	fmt.Printf("use time :%d ms\n", time.Since(now).Nanoseconds()/1e6)

	k, g := 0, 0
	for e1 := l.Front(); e1 != nil; e1 = e1.Next() {
		//t.Logf("id is %d\n", e1.Value.(goid.SFId))
		for e2 := e1.Next(); e2 != nil; e2 = e2.Next() {
			if e1.Value.(idUtil.SFId) == e2.Value.(idUtil.SFId) {
				for k, v := range e1.Value.(idUtil.SFId).Parse() {
					t.Logf("k :%s v: %d \n", k, v)
				}
				t.Fatalf("Same id! %d %d at %d %d\n", e1.Value.(idUtil.SFId), e2.Value.(idUtil.SFId), k, g)
			}
			g++
		}
		k++
		g = k + 1
	}
}
