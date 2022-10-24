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

type RecycleMap[K comparable, V any] interface {
	//设置一个值，含过期时间
	Set(key K, value V, expireIn time.Duration) error

	//根据key获取value
	Get(key K) V

	Delete(key K)

	//根据key设置key过期时间
	SetExpire(key K, expireIn time.Duration) bool

	//获得key过期时间
	TTL(key K) time.Duration

	//获得总数
	Size() int64

	Close() error
}

type dataEntity[V any] struct {
	value      V
	expireTime time.Time
}

type DeleteNotifier[K comparable, V any] func(key K, value V)

type defaultRecycleMap[K comparable, V any] struct {
	purgeInterval time.Duration
	purgeNumber   int64

	notifier DeleteNotifier[K, V]
	db       map[K]*dataEntity[V]
	stop     chan struct{}
	lock     sync.Locker
}

type Opt[K comparable, V any] func(*defaultRecycleMap[K, V])

func New[K comparable, V any](opts ...Opt[K, V]) RecycleMap[K, V] {
	ret := &defaultRecycleMap[K, V]{
		purgeInterval: DefaultPurgeInterval,
		purgeNumber:   DefaultPurgeNumberPerTime,
		db:            map[K]*dataEntity[V]{},
		stop:          make(chan struct{}),
		lock:          &sync.Mutex{},
	}
	for _, opt := range opts {
		opt(ret)
	}

	ret.run()
	return ret
}

func (dm *defaultRecycleMap[K, V]) purge() {
	dm.lock.Lock()
	defer dm.lock.Unlock()

	now := time.Now()
	var i int64 = 0
	for k, v := range dm.db {
		if v.expireTime.IsZero() {
			continue
		}
		if !v.expireTime.After(now) {
			dm.innerDelete(k, v.value)
			i++
			if i >= dm.purgeNumber {
				return
			}
		}
	}
}

// 初始化并开启回收线程，必须调用
func (dm *defaultRecycleMap[K, V]) run() {
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

// 关闭
func (dm *defaultRecycleMap[K, V]) Close() error {
	close(dm.stop)
	return nil
}

// 设置一个值，含过期时间
// 如果expireIn设置为-1，则永不过期
func (dm *defaultRecycleMap[K, V]) Set(key K, value V, expireIn time.Duration) error {
	dm.lock.Lock()
	defer dm.lock.Unlock()

	if expireIn >= 0 {
		dm.db[key] = &dataEntity[V]{value: value, expireTime: time.Now().Add(expireIn)}
	} else {
		dm.db[key] = &dataEntity[V]{value: value}
	}

	return nil
}

// 根据key获取value
func (dm *defaultRecycleMap[K, V]) Get(key K) V {
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
			dm.innerDelete(key, v.value)
			var v V
			return v
		}
		return v.value
	} else {
		var v V
		return v
	}
}

// 获得总数
func (dm *defaultRecycleMap[K, V]) Size() int64 {
	dm.lock.Lock()
	dm.lock.Unlock()

	var size int64 = 0
	now := time.Now()
	for k, v := range dm.db {
		if v.expireTime.IsZero() {
			size++
			continue
		}
		if !v.expireTime.After(now) {
			dm.innerDelete(k, v.value)
		} else {
			size++
		}
	}
	return size
}

// 删除key
func (dm *defaultRecycleMap[K, V]) Delete(key K) {
	dm.lock.Lock()
	defer dm.lock.Unlock()

	if v, ok := dm.db[key]; ok {
		dm.innerDelete(key, v.value)
	}
}

func (dm *defaultRecycleMap[K, V]) innerDelete(key K, value V) {
	delete(dm.db, key)
	dm.notifyDelete(key, value)
}

func (dm *defaultRecycleMap[K, V]) notifyDelete(key K, value V) {
	if dm.notifier != nil {
		dm.notifier(key, value)
	}
}

// 根据key设置key过期时间
func (dm *defaultRecycleMap[K, V]) SetExpire(key K, expireIn time.Duration) bool {
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

// 获得key过期时间, 如果不存在返回-2，如果永不过期返回-1
func (dm *defaultRecycleMap[K, V]) TTL(key K) time.Duration {
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
func OptSetPurgeInterval[K comparable, V any](interval time.Duration) Opt[K, V] {
	return func(recycleMap *defaultRecycleMap[K, V]) {
		recycleMap.purgeInterval = interval
	}
}

// 配置每次清理的清理数量，如果超出则放到下次清理(默认全部清理)
func OptSetPurgeNumberPerTime[K comparable, V any](number int64) Opt[K, V] {
	return func(recycleMap *defaultRecycleMap[K, V]) {
		recycleMap.purgeNumber = number
	}
}

// 配置清理回调函数
func OptSetDeleteNotifier[K comparable, V any](notifier DeleteNotifier[K, V]) Opt[K, V] {
	return func(recycleMap *defaultRecycleMap[K, V]) {
		recycleMap.notifier = notifier
	}
}

// 配置锁
func OptSetLocker[K comparable, V any](locker sync.Locker) Opt[K, V] {
	return func(recycleMap *defaultRecycleMap[K, V]) {
		recycleMap.lock = locker
	}
}

// 开启事务
func (dm *defaultRecycleMap[K, V]) Multi() error {
	//dm.Lock.Lock()
	return nil
}

// 执行事务
func (dm *defaultRecycleMap[K, V]) Exec() error {
	//dm.Lock.Unlock()
	return nil
}
