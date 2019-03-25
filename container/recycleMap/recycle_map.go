/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package recycleMap

import (
    "time"
    "sync"
    "runtime"
)

type dataEntity struct {
    value      interface{}
    expireTime time.Duration
}

type RecycleMap struct {
    PurgeInterval time.Duration
    db            map[interface{}]dataEntity
    stop          chan bool
    Lock          sync.Locker
}

func New() *RecycleMap {
    ret := &RecycleMap{
        db:   map[interface{}]dataEntity{},
        stop: make(chan bool),
        Lock: &sync.Mutex{},
    }

    return ret
}

func (dm *RecycleMap) purge() {
    dm.Lock.Lock()
    defer dm.Lock.Unlock()

    now := time.Duration(time.Now().UnixNano())
    for k, v := range dm.db {
        if v.expireTime <= now {
            delete(dm.db, k)
        }
    }
}

//初始化并开启回收线程，必须调用
func (dm *RecycleMap) Run() {
    if dm.PurgeInterval <= 0 {
        dm.PurgeInterval = 0
    }

    go func() {
        if dm.PurgeInterval > 0 {
            timer := time.NewTicker(dm.PurgeInterval)
            for {
                select {
                case <-dm.stop:
                    return
                case <-timer.C:
                    dm.purge()
                }
            }
        } else {
            for {
                select {
                case <-dm.stop:
                    return
                default:

                }
                dm.purge()

                runtime.Gosched()
            }
        }
    }()
}

//关闭
func (dm *RecycleMap) Close() {
    close(dm.stop)
}

//设置一个值，含过期时间
func (dm *RecycleMap) Set(key, value interface{}, expireIn time.Duration) {
    dm.Lock.Lock()
    defer dm.Lock.Unlock()

    time := expireIn + time.Duration(time.Now().UnixNano())
    dm.db[key] = dataEntity{value: value, expireTime: time}
}

//根据key获取value
func (dm *RecycleMap) Get(key interface{}) interface{} {
    dm.Lock.Lock()
    defer dm.Lock.Unlock()

    v, ok := dm.db[key]
    if ok {
        return v.value
    } else {
        return nil
    }
}

//删除key
func (dm *RecycleMap) Del(key interface{}) {
    dm.Lock.Lock()
    defer dm.Lock.Unlock()
    delete(dm.db, key)
}

//根据key设置key过期时间
func (dm *RecycleMap) SetExpire(key interface{}, expireIn time.Duration) bool {
    dm.Lock.Lock()
    defer dm.Lock.Unlock()

    v, ok := dm.db[key]
    if ok {
        v.expireTime = expireIn + time.Duration(time.Now().UnixNano())
        return true
    } else {
        return false
    }
}

//获得key过期时间
func (dm *RecycleMap) TTL(key interface{}) time.Duration {
    dm.Lock.Lock()
    defer dm.Lock.Unlock()

    v, ok := dm.db[key]
    if ok {
        return v.expireTime - time.Duration(time.Now().UnixNano())
    } else {
        return -1
    }
}

//开启事务
func (dm *RecycleMap) Multi() error {
    dm.Lock.Lock()
    return nil
}

//执行事务
func (dm *RecycleMap) Exec() error {
    dm.Lock.Unlock()
    return nil
}
