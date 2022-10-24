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
)

type RecycleMap[K comparable, V any] interface {
	// Set 设置一个值，含过期时间
	Set(key K, value V, expireIn time.Duration) error

	// Get 根据key获取value
	Get(key K) V

	// Keys 获得所有匹配的key
	Keys(pattern K) []K

	// Delete 删除key
	Delete(keys ...K) int64

	// SetExpire 根据key设置key过期时间
	SetExpire(key K, expireIn time.Duration) bool

	// TTL 获得key过期时间
	TTL(key K) time.Duration

	// Size 获得key总数
	Size() int64

	// Purge 回收过期key
	Purge()

	// Close 关闭并回收所有资源
	Close() error
}
