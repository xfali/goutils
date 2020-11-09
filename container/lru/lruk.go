// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package lru

import "container/list"

// history node
type hnode struct {
	hits int
	// for delete from map
	k interface{}
	// value
	v interface{}
}

func (n *hnode) clear() {
	n.hits = 0
	n.k = nil
	n.v = nil
}

type LRUK struct {
	m map[interface{}]*list.Element
	// history queue
	hQueue *list.List
	// cache queue
	cQueue *list.List

	k         int
	hQueueCap int
	cQueueCap int
}

func NewLruKCache(k, historyCapacity, cacheCapacity int) *LRUK {
	ret := &LRUK{
		m:         map[interface{}]*list.Element{},
		hQueue:    list.New(),
		cQueue:    list.New(),
		hQueueCap: historyCapacity,
		cQueueCap: cacheCapacity,
		k:         k,
	}
	return ret
}

func (m *LRUK) hit(key interface{}, hit bool) {

}

func (m *LRUK) insert(key interface{}, value interface{}) {
	if m.Size()+1 > m.hQueueCap {
		e := m.hQueue.Back()
		if e != nil {
			n := e.Value.(*hnode)
			m.hQueue.Remove(e)
			// not cache
			if n.hits < m.k {
				delete(m.m, e.Value.([2]interface{})[0])
			}
			n.clear()
		}
	}
	m.m[key] = m.hQueue.PushFront(&hnode{k: key, v: value})
}

// 向Map中添加一个元素
// Param：key 添加的对象key，value 添加的对象
func (m *LRUK) Put(key, value interface{}) {
	e, ok := m.m[key]
	if ok {
		e.Value.(*hnode).clear()
		e.Value = &hnode{k: key, v: value}
		m.hQueue.MoveToFront(e)
	} else {
		m.insert(key, value)
	}
}

// 获取key对应的元素
// Param：key 对象key
// Return： value：key对应的对象，loaded：成功获取返回true，不存在返回false
func (m *LRUK) Get(key interface{}) (value interface{}, loaded bool) {
	v, ok := m.m[key]
	if ok {
		m.hQueue.MoveToFront(v)
		m.hit(key, true)
		return v.Value.(*hnode).v, true
	} else {
		m.hit(key, false)
		return nil, false
	}
}

// 删除key对应的元素
// Param：key
func (m *LRUK) Delete(key interface{}) {
	v, ok := m.m[key]
	if ok {
		m.hQueue.Remove(v)
		delete(m.m, key)
	}
}

// 获得Map长度
// Return： 链表长度
func (m *LRUK) Size() int {
	return len(m.m)
}
