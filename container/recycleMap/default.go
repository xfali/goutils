/*
 * Copyright 2022 Xiongfa Li.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package recycleMap

import (
	"fmt"
	"github.com/xfali/goutils/v2/container/purger"
	"math"
	"reflect"
	"regexp"
	"sync"
	"time"
)

const (
	// DefaultPurgeInterval 默认清理时间间隔
	DefaultPurgeInterval = 250 * time.Millisecond

	// DefaultPurgeNumberPerTime 默认每次清理个数
	DefaultPurgeNumberPerTime = math.MaxInt64
)

type dataEntity[V any] struct {
	value      V
	expireTime time.Time
}

type DeleteNotifier[K comparable, V any] func(key K, value V)

type defaultRecycleMap[K comparable, V any] struct {
	purgeInterval time.Duration
	purgeNumber   int64
	manualPurge   bool

	matcher  MatchFunc[K]
	notifier DeleteNotifier[K, V]
	db       map[K]*dataEntity[V]
	lock     sync.Locker
}

type Opt[K comparable, V any] func(*defaultRecycleMap[K, V])
type MatchFunc[K comparable] func(pattern K) func(key K) (match bool)

func New[K comparable, V any](opts ...Opt[K, V]) RecycleMap[K, V] {
	ret := &defaultRecycleMap[K, V]{
		purgeInterval: DefaultPurgeInterval,
		purgeNumber:   DefaultPurgeNumberPerTime,
		db:            map[K]*dataEntity[V]{},
		matcher:       defaultMatch[K],
		lock:          &sync.Mutex{},
	}
	for _, opt := range opts {
		opt(ret)
	}

	if !ret.manualPurge {
		err := execInstance().AddPurger(ret, ret.purgeInterval)
		if err != nil {
			panic(err)
		}
	}

	return ret
}

func (dm *defaultRecycleMap[K, V]) Keys(pattern K) []K {
	dm.lock.Lock()
	defer dm.lock.Unlock()

	matcher := dm.matcher(pattern)
	ret := make([]K, 0, len(dm.db))
	for k := range dm.db {
		if matcher(k) {
			ret = append(ret, k)
		}
	}
	return ret
}

func (dm *defaultRecycleMap[K, V]) Purge() {
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

// 关闭
func (dm *defaultRecycleMap[K, V]) Close() error {
	dm.Purge()
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
func (dm *defaultRecycleMap[K, V]) Delete(keys ...K) int64 {
	dm.lock.Lock()
	defer dm.lock.Unlock()

	var total int64 = 0
	for _, key := range keys {
		if v, ok := dm.db[key]; ok {
			dm.innerDelete(key, v.value)
			total++
		}
	}
	return total
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

func AllMatch[K comparable](pattern K) func(key K) (match bool) {
	return func(key K) (match bool) {
		return true
	}
}

func defaultMatch[K comparable](pattern K) func(key K) (match bool) {
	return func(key K) (match bool) {
		defer func(ret *bool) {
			if o := recover(); o != nil {
				*ret = pattern == key
			}
		}(&match)
		if v := reflect.ValueOf(pattern); v.IsZero() {
			return true
		}
		return pattern == key
	}
}

func RegexpMatcher[K comparable](converter func(v K) string) func(pattern K) func(key K) (match bool) {
	if converter == nil {
		converter = func(v K) string {
			return fmt.Sprintf("%v", v)
		}
	}
	return func(pattern K) func(key K) (match bool) {
		reg, err := regexp.Compile(converter(pattern))
		if err != nil {
			return func(key K) (match bool) {
				return false
			}
		}
		return func(key K) (match bool) {
			return reg.MatchString(converter(key))
		}
	}
}

// OptManualPurge 配置手工清理，注意过期key将不再自动清理，但是还是会触发被动清理
func OptManualPurge[K comparable, V any]() Opt[K, V] {
	return func(recycleMap *defaultRecycleMap[K, V]) {
		recycleMap.manualPurge = true
	}
}

// OptAutoPurge 配置清理时间间隔
// 设置时间间隔越短清理得越及时，但是消耗更多CPU
// 设置时间间隔越长内存消耗越多。
func OptAutoPurge[K comparable, V any](interval time.Duration, executor purger.PurgeExecutor) Opt[K, V] {
	return func(recycleMap *defaultRecycleMap[K, V]) {
		if executor == nil {
			executor = execInstance()
		}
		recycleMap.manualPurge = true
		if err := executor.AddPurger(recycleMap, interval); err != nil {
			panic(err)
		}
	}
}

// OptSetPurgeInterval 配置清理时间间隔（默认50ms）
// 设置时间间隔越短清理得越及时，但是消耗更多CPU
// 设置时间间隔越长内存消耗越多。
func OptSetPurgeInterval[K comparable, V any](interval time.Duration) Opt[K, V] {
	return func(recycleMap *defaultRecycleMap[K, V]) {
		recycleMap.purgeInterval = interval
	}
}

// OptSetPurgeNumberPerTime 配置每次清理的清理数量，如果超出则放到下次清理(默认全部清理)
func OptSetPurgeNumberPerTime[K comparable, V any](number int64) Opt[K, V] {
	return func(recycleMap *defaultRecycleMap[K, V]) {
		recycleMap.purgeNumber = number
	}
}

// OptSetDeleteNotifier 配置清理回调函数
func OptSetDeleteNotifier[K comparable, V any](notifier DeleteNotifier[K, V]) Opt[K, V] {
	return func(recycleMap *defaultRecycleMap[K, V]) {
		recycleMap.notifier = notifier
	}
}

// OptSetLocker 配置锁
func OptSetLocker[K comparable, V any](locker sync.Locker) Opt[K, V] {
	return func(recycleMap *defaultRecycleMap[K, V]) {
		recycleMap.lock = locker
	}
}

// OptSetMatcher 配置锁
func OptSetMatcher[K comparable, V any](matcher MatchFunc[K]) Opt[K, V] {
	return func(recycleMap *defaultRecycleMap[K, V]) {
		recycleMap.matcher = matcher
	}
}

// Multi 开启事务
func (dm *defaultRecycleMap[K, V]) Multi() error {
	//dm.Lock.Lock()
	return nil
}

// Exec 执行事务
func (dm *defaultRecycleMap[K, V]) Exec() error {
	//dm.Lock.Unlock()
	return nil
}

var (
	gExecutor   purger.PurgeExecutor
	execCreator = func() purger.PurgeExecutor {
		return purger.New()
	}
	initOnce sync.Once
)

func execInstance() purger.PurgeExecutor {
	initOnce.Do(func() {
		gExecutor = execCreator()
	})
	return gExecutor
}

// CloseGlobalPurgeExecutor 关闭全局自动回收器
func CloseGlobalPurgeExecutor() error {
	return execInstance().Close()
}

type Config struct {
	// 自动回收执行器创建方法
	PurgeExecutorCreator func() purger.PurgeExecutor
}

// Init 初始化函数，必须在package引用时调用
func Init(conf Config) error {
	if conf.PurgeExecutorCreator != nil {
		execCreator = conf.PurgeExecutorCreator
	}
	return nil
}
