// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package lru

type LRU interface {
	Get(key interface{}) (value interface{})
	Put(key, value interface{})
}
