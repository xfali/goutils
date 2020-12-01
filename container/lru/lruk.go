// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package lru

// history node
type hnode struct {
	hits int
	// for delete from map
	k interface{}
	// value
	v interface{}

	he *QueueElem
	ce *QueueElem
}

func (n *hnode) clear() {
	n.hits = 0
	n.k = nil
	n.v = nil

	n.he = nil
	n.ce = nil
}

type LRUK struct {
	m map[interface{}]*hnode
	// history queue
	hQueue *LruQueue
	// cache queue
	cQueue *LruQueue

	purgeFunc func()
}

func NewLruKCache(k, historyCapacity, cacheCapacity int) *LRUK {
	ret := &LRUK{
		m: map[interface{}]*hnode{},
	}
	ret.cQueue = NewLruQueue(cacheCapacity)
	cl := &cacheListener{lru: ret}
	ret.cQueue.AddListener(cl)
	ret.hQueue = NewLruQueue(historyCapacity)
	hl := &historyListener{lru: ret, history: ret.hQueue, cache: ret.cQueue, k: k}
	ret.hQueue.AddListener(hl)
	ret.purgeFunc = func() {
		cl.lru = nil
		hl.lru = nil
		hl.history = nil
		hl.cache = nil
	}

	return ret
}

func (m *LRUK) hit(key interface{}, hit bool) {

}

func (m *LRUK) Purge() {
	m.cQueue.listeners = nil
	m.cQueue.list = nil
	m.cQueue = nil

	m.hQueue.listeners = nil
	m.hQueue.list = nil
	m.hQueue = nil
	m.purgeFunc()
}

// 向Map中添加一个元素
// Param：key 添加的对象key，value 添加的对象
func (m *LRUK) Put(key, value interface{}) {
	e, ok := m.m[key]
	if ok {
		n := e
		if n.he != nil {
			m.hQueue.Delete(n.he)
		}
		if n.ce != nil {
			m.cQueue.Delete(n.ce)
		}
		n.clear()
	}
	e = &hnode{k: key, v: value}
	elem := m.hQueue.Insert(e)
	e.he = elem
	m.m[key] = e
}

// 获取key对应的元素
// Param：key 对象key
// Return： value：key对应的对象，loaded：成功获取返回true，不存在返回false
func (m *LRUK) Get(key interface{}) (value interface{}, loaded bool) {
	v, ok := m.m[key]
	if ok {
		if m.checkAndDelete(v) {
			return nil, false
		}
		if v.he != nil {
			m.hQueue.Touch(v.he)
		} else {
			elem := m.hQueue.Insert(v)
			v.he = elem
		}
		m.hit(key, true)
		return v.v, true
	} else {
		m.hit(key, false)
		return nil, false
	}
}

func (m *LRUK) checkAndDelete(v *hnode) bool {
	if v.he == nil && v.ce == nil {
		delete(m.m, v.k)
		return true
	}
	return false
}

func (m *LRUK) delete(k interface{}, v *hnode) {
	if v.he != nil {
		m.hQueue.Delete(v.he)
	}
	if v.ce != nil {
		m.cQueue.Delete(v.ce)
	}
	delete(m.m, k)
}

// 删除key对应的元素
// Param：key
func (m *LRUK) Delete(key interface{}) {
	v, ok := m.m[key]
	if ok {
		m.delete(key, v)
	}
}

// 获得Map长度
// Return： 链表长度
func (m *LRUK) Size() int {
	return len(m.m)
}

type historyListener struct {
	lru     *LRUK
	history *LruQueue
	cache   *LruQueue
	k       int
}

func (m *historyListener) PostTouch(v interface{}) {
	n := v.(*hnode)
	n.hits++
	//move to cache
	if n.hits >= m.k {
		e := m.cache.Insert(n)
		n.ce = e
		m.history.Delete(n.he)
		n.he = nil
	}
}

func (m *historyListener) PostInsert(v interface{}) {
}

func (m *historyListener) PostDelete(v interface{}) {
	n := v.(*hnode)
	n.he = nil
	m.lru.checkAndDelete(n)
}

type cacheListener struct {
	lru *LRUK
}

func (m *cacheListener) PostTouch(v interface{}) {
}

func (m *cacheListener) PostInsert(v interface{}) {
}

func (m *cacheListener) PostDelete(v interface{}) {
	n := v.(*hnode)
	n.ce = nil
	m.lru.checkAndDelete(n)
}
