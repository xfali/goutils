/**
 * Copyright (C) 2018, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @date 2018/7/23 
 * @time 16:15
 * @version V1.0
 * Description: 
 */

package goutils

import (
    "container/list"
    "sync"
)

type BlockQueue struct {
    list *list.List
    maxSize int
    lock *sync.Mutex
    cond *sync.Cond
}

type OnEnqueque func(data interface{})(bool)

func NewBlockQueue(maxSize int) (*BlockQueue) {
    lock := &sync.Mutex{}
    return &BlockQueue{
        list: list.New(),
        maxSize: maxSize,
        lock: lock,
        cond: sync.NewCond(lock),
    }
}

func (bq *BlockQueue)Enqueue(data interface{}) {
    bq.lock.Lock()
    defer bq.lock.Unlock()

    for bq.list.Len() == bq.maxSize {
        bq.cond.Wait()
    }

    if bq.list.Len() == 0 {
        defer bq.cond.Broadcast()
    }

    bq.list.PushBack(data)
}

func (bq *BlockQueue)TryEnqueue(data interface{}) (bool) {
    bq.lock.Lock()
    defer bq.lock.Unlock()

    for bq.list.Len() == bq.maxSize {
        return false
    }

    if bq.list.Len() == 0 {
        defer bq.cond.Broadcast()
    }

    bq.list.PushBack(data)
    return true
}

func (bq *BlockQueue) First() interface{} {
    bq.lock.Lock()
    defer bq.lock.Unlock()
    elm := bq.list.Front()
    return elm.Value
}

func (bq *BlockQueue) WaitOne(onFunc OnEnqueque) {
    if onFunc == nil {
        return
    }

    bq.lock.Lock()
    defer bq.lock.Unlock()

    for bq.list.Len() == 0 {
        bq.cond.Wait()
    }

    if bq.list.Len() == bq.maxSize {
        defer bq.cond.Broadcast()
    }
    elm := bq.list.Front()
    if onFunc(elm.Value) {
        bq.list.Remove(elm)
    }
}

func (bq *BlockQueue) Dequeue() interface{} {
    bq.lock.Lock()
    defer bq.lock.Unlock()

    for bq.list.Len() == 0 {
        bq.cond.Wait()
    }

    if bq.list.Len() == bq.maxSize {
        defer bq.cond.Broadcast()
    }

    elm := bq.list.Front()
    return bq.list.Remove(elm)
}

