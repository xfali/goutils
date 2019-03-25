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

    ret.Init()

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

func (dm *RecycleMap) Init() {
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

func (dm *RecycleMap) Close() {
    close(dm.stop)
}

func (dm *RecycleMap) Set(key, value interface{}, expireIn time.Duration) {
    dm.Lock.Lock()
    defer dm.Lock.Unlock()

    time := expireIn + time.Duration(time.Now().UnixNano())
    dm.db[key] = dataEntity{value: value, expireTime: time}
}

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

func (dm *RecycleMap) Del(key interface{}) {
    dm.Lock.Lock()
    defer dm.Lock.Unlock()
    delete(dm.db, key)
}

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

func (dm *RecycleMap) TTL(key string) time.Duration {
    dm.Lock.Lock()
    defer dm.Lock.Unlock()

    v, ok := dm.db[key]
    if ok {
        return v.expireTime - time.Duration(time.Now().UnixNano())
    } else {
        return -1
    }
}

func (dm *RecycleMap) Multi() error {
    dm.Lock.Lock()
    return nil
}

func (dm *RecycleMap) Exec() error {
    dm.Lock.Unlock()
    return nil
}
