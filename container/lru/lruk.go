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
}

func (n *hnode) clear() {
	n.hits = 0
	n.k = nil
	n.v = nil
}

type LRUK struct {
	m map[interface{}]*QueueElem
	// history queue
	hQueue *LruQueue
	// cache queue
	cQueue *SimpleLru
}

func NewLruKCache(k, historyCapacity, cacheCapacity int) *LRUK {
	ret := &LRUK{
		m: map[interface{}]*QueueElem{},
	}
	ret.hQueue = NewLruQueue(historyCapacity)
	ret.hQueue.AddListener(&historyListener{cache: ret, k: k})
	ret.cQueue = NewLruCache(cacheCapacity)
	return ret
}

func (m *LRUK) hit(key interface{}, hit bool) {

}

// 向Map中添加一个元素
// Param：key 添加的对象key，value 添加的对象
func (m *LRUK) Put(key, value interface{}) {
	e, ok := m.m[key]
	if ok {
		e.Value.(*hnode).clear()
		e.Value = &hnode{k: key, v: value}
		m.hQueue.Touch(e)
	} else {
		elem := m.hQueue.Insert(&hnode{k: key, v: value})
		m.m[key] = elem
	}
}

// 获取key对应的元素
// Param：key 对象key
// Return： value：key对应的对象，loaded：成功获取返回true，不存在返回false
func (m *LRUK) Get(key interface{}) (value interface{}, loaded bool) {
	v, ok := m.m[key]
	if ok {
		m.hQueue.Touch(v)
		m.hit(key, true)
		return v.Value.(*hnode).v, true
	} else {
		m.hit(key, false)
		return nil, false
	}
}

func (m *LRUK) delete(k interface{}, v *QueueElem) {
	m.cQueue.Delete(k)
	m.hQueue.Delete(v)
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
	cache *LRUK
	k     int
}

func (m *historyListener) PostTouch(v interface{}) {
	n := v.(hnode)
	n.hits++
	if n.hits > m.k {
		m.cache.cQueue.Put(n.k, nil)
	}
}

func (m *historyListener) PostInsert(v interface{}) {
}

func (m *historyListener) PostDelete(v interface{}) {
	n := v.(*hnode)
	_, ok := m.cache.cQueue.Get(n.k)
	if !ok {
		delete(m.cache.m, n.k)
	}
	n.clear()
}

type cacheListener struct {
	cache *LRUK
}

func (m *cacheListener) PostTouch(v interface{}) {
}

func (m *cacheListener) PostInsert(v interface{}) {
}

func (m *cacheListener) PostDelete(v interface{}) {
	n := v.(*hnode)
	_, ok := m.cache.cQueue.Get(n.k)
	if !ok {
		delete(m.cache.m, n.k)
	}
	n.clear()
}
