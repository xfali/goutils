// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package lru

import (
	"github.com/xfali/goutils/v2/container/xmap"
)

type LRU interface {
	xmap.IMap

	Purge()
}

type SimpleLru struct {
	m     map[interface{}]*QueueElem
	queue *LruQueue

	cap int
}

func NewLruCache(capacity int) *SimpleLru {
	ret := &SimpleLru{
		m:   map[interface{}]*QueueElem{},
		cap: capacity,
	}
	ret.queue = NewLruQueue(capacity)
	ret.queue.AddListener(ret)
	return ret
}

func (m *SimpleLru) PostTouch(v interface{}) {
}

func (m *SimpleLru) PostInsert(v interface{}) {
}

func (m *SimpleLru) PostDelete(v interface{}) {
	key := v.([2]interface{})[0]
	delete(m.m, key)
}

func (m *SimpleLru) hit(key interface{}, hit bool) {

}

func (m *SimpleLru) Purge() {
	m.queue.listeners = nil
	m.queue.list = nil
	m.queue = nil
}

// 向Map中添加一个元素
// Param：key 添加的对象key，value 添加的对象
func (m *SimpleLru) Put(key, value interface{}) {
	e, ok := m.m[key]
	if ok {
		e.Value = [2]interface{}{key, value}
		m.queue.Touch(e)
	} else {
		elem := m.queue.Insert([2]interface{}{key, value})
		m.m[key] = elem
	}
}

// 获取key对应的元素
// Param：key 对象key
// Return： value：key对应的对象，loaded：成功获取返回true，不存在返回false
func (m *SimpleLru) Get(key interface{}) (value interface{}, loaded bool) {
	v, ok := m.m[key]
	if ok {
		m.queue.Touch(v)
		// 命中
		m.hit(key, true)
		return v.Value.([2]interface{})[1], true
	} else {
		// 未命中
		m.hit(key, false)
		return nil, false
	}
}

// 删除key对应的元素
// Param：key
func (m *SimpleLru) Delete(key interface{}) {
	v, ok := m.m[key]
	if ok {
		m.queue.Delete(v)
	}
}

// 获得Map长度
// Return： 链表长度
func (m *SimpleLru) Size() int {
	return len(m.m)
}
