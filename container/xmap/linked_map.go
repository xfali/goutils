// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package xmap

import "container/list"

type node struct {
	prev *node
	next *node
	v    interface{}
}

func (n *node) init() {
	n.prev = nil
	n.next = nil
	n.v = nil
}

type LinkedMap struct {
	head node
	m    map[interface{}]*node
}

func NewLinkedMap() *LinkedMap {
	ret := &LinkedMap{
		m: make(map[interface{}]*node),
	}
	ret.head.init()
	return ret
}

func (m *LinkedMap) init() {
	if m.head.next == nil {
		m.head.next = &m.head
		m.head.prev = &m.head
	}
}

func (m *LinkedMap) insert(v interface{}, at *node) *node {
	e := &node{v: v}
	n := at.next
	at.next = e
	e.prev = at
	e.next = n
	n.prev = e

	return e
}

// 向Map中添加一个元素
// Param：key 添加的对象key，value 添加的对象
// Return： 成功返回true，失败返回false
func (m *LinkedMap) Put(key, value interface{}) {
	if e, ok := m.m[key]; ok {
		e.v = [2]interface{}{key, value}
		return
	}

	m.init()
	n := m.insert([2]interface{}{key, value}, m.head.prev)
	m.m[key] = n
}

// 尝试向Map中添加一个元素，如果已存在该元素则直接返回已存在元素不进行添加
// Param：key 添加的对象key，value 添加的对象
// Return： actual 如果key已存在对应元素，则返回该元素，否则返回新添加的元素。 loaded：已存在返回true，否则返回false
func (m *LinkedMap) GetOrPut(key, value interface{}) (actual interface{}, loaded bool) {
	o, ok := m.m[key]
	if ok {
		return o.v.([2]interface{})[1], true
	}

	m.init()
	n := m.insert([2]interface{}{key, value}, m.head.prev)
	m.m[key] = n
	return value, false
}

// 获取key对应的元素
// Param：key 对象key
// Return： value：key对应的对象，loaded：成功获取返回true，不存在返回false
func (m *LinkedMap) Get(key interface{}) (value interface{}, loaded bool) {
	o, ok := m.m[key]
	if ok {
		return o.v.([2]interface{})[1], ok
	}
	return nil, false
}

// 删除key对应的元素
// Param：key
func (m *LinkedMap) Delete(key interface{}) {
	if n, ok := m.m[key]; ok {
		n.prev.next = n.next
		n.next.prev = n.prev
		n.init()
		delete(m.m, key)
	}
}

// 获得Map长度
// Return： 链表长度
func (m *LinkedMap) Size() int {
	return len(m.m)
}

// 轮询Map O(N)
// Param：接受轮询的函数，返回true继续轮询，返回false终止轮询
func (m *LinkedMap) Foreach(f func(key interface{}, value interface{}) bool) {
	for e := m.head.next; e != nil && e != &m.head; e = e.next {
		kv := e.v.([2]interface{})
		if !f(kv[0], kv[1]) {
			break
		}
	}
}

// 查询Map中是否存在参数对象
// Param：查询的对象
// Return：存在返回true，不存在返回false
func (m *LinkedMap) Find(key interface{}) bool {
	_, ok := m.m[key]
	return ok
}

type SimpleLinkedMap struct {
	l *list.List
	m map[interface{}]*list.Element
}

func NewSimpleLinkedMap() *SimpleLinkedMap {
	return &SimpleLinkedMap{list.New(), make(map[interface{}]*list.Element)}
}

// 向Map中添加一个元素
// Param：key 添加的对象key，value 添加的对象
// Return： 成功返回true，失败返回false
func (m *SimpleLinkedMap) Put(key, value interface{}) {
	if e, ok := m.m[key]; ok {
		e.Value = [2]interface{}{key, value}
	}
	m.m[key] = m.l.PushBack([2]interface{}{key, value})
}

// 尝试向Map中添加一个元素，如果已存在该元素则直接返回已存在元素不进行添加
// Param：key 添加的对象key，value 添加的对象
// Return： actual 如果key已存在对应元素，则返回该元素，否则返回新添加的元素。 loaded：已存在返回true，否则返回false
func (m *SimpleLinkedMap) GetOrPut(key, value interface{}) (actual interface{}, loaded bool) {
	o, ok := m.m[key]
	if ok {
		return o.Value.([2]interface{})[1], true
	}
	m.m[key] = m.l.PushBack([2]interface{}{key, value})
	return value, false
}

// 获取key对应的元素
// Param：key 对象key
// Return： value：key对应的对象，loaded：成功获取返回true，不存在返回false
func (m *SimpleLinkedMap) Get(key interface{}) (value interface{}, loaded bool) {
	o, ok := m.m[key]
	if ok {
		return o.Value.([2]interface{})[1], ok
	}
	return nil, false
}

// 删除key对应的元素
// Param：key
func (m *SimpleLinkedMap) Delete(key interface{}) {
	if e, ok := m.m[key]; ok {
		m.l.Remove(e)
		delete(m.m, key)
	}
}

// 获得Map长度
// Return： 链表长度
func (m *SimpleLinkedMap) Size() int {
	return len(m.m)
}

// 轮询Map O(N)
// Param：接受轮询的函数，返回true继续轮询，返回false终止轮询
func (m *SimpleLinkedMap) Foreach(f func(key interface{}, value interface{}) bool) {
	for e := m.l.Front(); e != nil; e = e.Next() {
		kv := e.Value.([2]interface{})
		if !f(kv[0], kv[1]) {
			break
		}
	}
}

// 查询Map中是否存在参数对象
// Param：查询的对象
// Return：存在返回true，不存在返回false
func (m *SimpleLinkedMap) Find(key interface{}) bool {
	_, ok := m.m[key]
	return ok
}
