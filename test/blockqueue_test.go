/**
 * Copyright (C) 2018, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @date 2018/7/23
 * @time 16:32
 * @version V1.0
 * Description:
 */

package test

import (
	"fmt"
	"github.com/xfali/goutils/container"
	"testing"
	"time"
)

func TestBlockQueue1(t *testing.T) {
	bq := container.NewBlockQueue(2)
	stop := false

	go func() {
		for !stop {
			time.Sleep(time.Second)
			fmt.Println(bq.Dequeue().(int))
		}
	}()

	for i := 0; i < 10; i++ {
		bq.Enqueue(i)
	}

	select {
	case <-time.After(15 * time.Second):
		return
	}
	stop = true
}

func TestBlockQueue2(t *testing.T) {
	bq := container.NewBlockQueue(2)
	stop := false

	go func() {
		for !stop {
			fmt.Println(bq.Dequeue().(int))
		}
	}()

	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)
		bq.Enqueue(i)
	}

	select {
	case <-time.After(15 * time.Second):
		return
	}
	stop = true
}

func TestBlockQueue3(t *testing.T) {
	bq := container.NewBlockQueue(2)
	stop := false

	go func() {
		test := true
		for !stop {
			bq.WaitOne(func(data interface{}) bool {
				fmt.Println(data)

				if test {
					test = !test
					return false
				} else {
					test = !test
					return true
				}
			})
		}
	}()

	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)
		bq.Enqueue(i)
	}

	select {
	case <-time.After(15 * time.Second):
		return
	}
	stop = true
}
