/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description:
 */

package recycleMap

import (
	"math"
	"runtime"
	"sync"
	"time"
)

const (
	// 默认清理时间间隔
	DefaultPurgeInterval = 50 * time.Millisecond
	// 默认每次清理个数
	DefaultPurgeNumberPerTime = math.MaxInt64
)

type RecycleMap interface {
	//设置一个值，含过期时间
	Set(key, value interface{}, expireIn time.Duration) error

	//根据key获取value
	Get(key interface{}) interface{}

	Delete(key interface{})

	//根据key设置key过期时间
	SetExpire(key interface{}, expireIn time.Duration) bool

	//获得key过期时间
	TTL(key interface{}) time.Duration

	//获得总数
	Size() int64

	Close() error
}

type dataEntity struct {
	value      interface{}
	expireTime time.Time
}

type defaultRecycleMap struct {
	purgeInterval time.Duration
	purgeNumber   int64

	db   map[interface{}]*dataEntity
	stop chan bool
	lock sync.Locker
}

type Opt func(*defaultRecycleMap)

func New(opts ...Opt) RecycleMap {
	ret := &defaultRecycleMap{
		purgeInterval: DefaultPurgeInterval,
		purgeNumber:   DefaultPurgeNumberPerTime,
		db:            map[interface{}]*dataEntity{},
		stop:          make(chan bool),
		lock:          &sync.Mutex{},
	}
	for _, opt := range opts {
		opt(ret)
	}

	ret.run()
	return ret
}

func (dm *defaultRecycleMap) purge() {
	dm.lock.Lock()
	defer dm.lock.Unlock()

	now := time.Now()
	var i int64 = 0
	for k, v := range dm.db {
		if v.expireTime.IsZero() {
			continue
		}
		if !v.expireTime.After(now) {
			delete(dm.db, k)
			i++
			if i >= dm.purgeNumber {
				return
			}
		}
	}
}

//初始化并开启回收线程，必须调用
func (dm *defaultRecycleMap) run() {
	if dm.purgeInterval <= 0 {
		dm.purgeInterval = 0
	}

	go func() {
		if dm.purgeInterval > 0 {
			timer := time.NewTicker(dm.purgeInterval)
			defer timer.Stop()
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
func (dm *defaultRecycleMap) Close() error {
	close(dm.stop)
	return nil
}

// 设置一个值，含过期时间
// 如果expireIn设置为-1，则永不过期
func (dm *defaultRecycleMap) Set(key, value interface{}, expireIn time.Duration) error {
	dm.lock.Lock()
	defer dm.lock.Unlock()

	if expireIn >= 0 {
		dm.db[key] = &dataEntity{value: value, expireTime: time.Now().Add(expireIn)}
	} else {
		dm.db[key] = &dataEntity{value: value}
	}

	return nil
}

//根据key获取value
func (dm *defaultRecycleMap) Get(key interface{}) interface{} {
	dm.lock.Lock()
	defer dm.lock.Unlock()

	v, ok := dm.db[key]
	if ok {
		// expired
		if v.expireTime.IsZero() {
			return v.value
		}
		now := time.Now()
		if !v.expireTime.After(now) {
			return nil
		}
		return v.value
	} else {
		return nil
	}
}

//获得总数
func (dm *defaultRecycleMap) Size() int64 {
	dm.lock.Lock()
	dm.lock.Unlock()

	var size int64 = 0
	now := time.Now()
	for _, v := range dm.db {
		if v.expireTime.IsZero() {
			size++
			continue
		}
		if !v.expireTime.After(now) {
		} else {
			size++
		}
	}
	return size
}

//删除key
func (dm *defaultRecycleMap) Delete(key interface{}) {
	dm.lock.Lock()
	defer dm.lock.Unlock()
	delete(dm.db, key)
}

//根据key设置key过期时间
func (dm *defaultRecycleMap) SetExpire(key interface{}, expireIn time.Duration) bool {
	dm.lock.Lock()
	defer dm.lock.Unlock()

	v, ok := dm.db[key]
	if ok {
		if expireIn >= 0 {
			v.expireTime = time.Now().Add(expireIn)
		} else {
			v.expireTime = time.Time{}
		}
		return true
	} else {
		return false
	}
}

//获得key过期时间, 如果不存在返回-2，如果永不过期返回-1
func (dm *defaultRecycleMap) TTL(key interface{}) time.Duration {
	dm.lock.Lock()
	defer dm.lock.Unlock()

	v, ok := dm.db[key]
	if ok {
		if v.expireTime.IsZero() {
			return -1
		}
		return v.expireTime.Sub(time.Now())
	} else {
		return -2
	}
}

// 配置清理时间间隔（默认50ms）
// 设置时间间隔越短清理得越及时，但是消耗更多CPU
// 设置时间间隔越长内存消耗越多。
func OptSetPurgeInterval(interval time.Duration) Opt {
	return func(recycleMap *defaultRecycleMap) {
		recycleMap.purgeInterval = interval
	}
}

// 配置每次清理的清理数量，如果超出则放到下次清理(默认全部清理)
func OptSetPurgeNumberPerTime(number int64) Opt {
	return func(recycleMap *defaultRecycleMap) {
		recycleMap.purgeNumber = number
	}
}

//开启事务
func (dm *defaultRecycleMap) Multi() error {
	//dm.Lock.Lock()
	return nil
}

//执行事务
func (dm *defaultRecycleMap) Exec() error {
	//dm.Lock.Unlock()
	return nil
}
