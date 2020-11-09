// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package lru

import (
	"container/list"
	"github.com/xfali/goutils/container/xmap"
)

type LRU interface {
	xmap.IMap
}

type SimpleLru struct {
	m map[interface{}]*list.Element
	l *list.List

	cap int
}

func NewLruCache(capacity int) *SimpleLru {
	ret := &SimpleLru{
		m:   map[interface{}]*list.Element{},
		l:   list.New(),
		cap: capacity,
	}
	return ret
}

func (m *SimpleLru) hit(key interface{}, hit bool) {

}

func (m *SimpleLru) insert(key interface{}, value interface{}) {
	if m.Size()+1 > m.cap {
		e := m.l.Back()
		if e != nil {
			m.l.Remove(e)
			delete(m.m, e.Value.([2]interface{})[0])
		}
	}
	m.m[key] = m.l.PushFront([2]interface{}{key, value})
}

// 向Map中添加一个元素
// Param：key 添加的对象key，value 添加的对象
func (m *SimpleLru) Put(key, value interface{}) {
	e, ok := m.m[key]
	if ok {
		e.Value = [2]interface{}{key, value}
		m.l.MoveToFront(e)
	} else {
		m.insert(key, value)
	}
}

// 获取key对应的元素
// Param：key 对象key
// Return： value：key对应的对象，loaded：成功获取返回true，不存在返回false
func (m *SimpleLru) Get(key interface{}) (value interface{}, loaded bool) {
	v, ok := m.m[key]
	if ok {
		m.l.MoveToFront(v)
		m.hit(key, true)
		return v.Value.([2]interface{})[1], true
	} else {
		m.hit(key, false)
		return nil, false
	}
}

// 删除key对应的元素
// Param：key
func (m *SimpleLru) Delete(key interface{}) {
	v, ok := m.m[key]
	if ok {
		m.l.Remove(v)
		delete(m.m, key)
	}
}

// 获得Map长度
// Return： 链表长度
func (m *SimpleLru) Size() int {
	return m.l.Len()
}
