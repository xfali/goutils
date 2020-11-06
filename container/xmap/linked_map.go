// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package xmap

import "container/list"

type LinkedMap struct {
	l *list.List
	m map[interface{}]*list.Element
}

func NewLinkedMap() *LinkedMap {
	return &LinkedMap{list.New(), make(map[interface{}]*list.Element)}
}

// 向Map中添加一个元素
// Param：key 添加的对象key，value 添加的对象
// Return： 成功返回true，失败返回false
func (m LinkedMap) Put(key, value interface{}) {
	if e, ok := m.m[key]; ok {
		e.Value = [2]interface{}{key, value}
	}
	m.m[key] = m.l.PushBack([2]interface{}{key, value})
}

// 尝试向Map中添加一个元素，如果已存在该元素则直接返回已存在元素不进行添加
// Param：key 添加的对象key，value 添加的对象
// Return： actual 如果key已存在对应元素，则返回该元素，否则返回新添加的元素。 loaded：已存在返回true，否则返回false
func (m LinkedMap) GetOrPut(key, value interface{}) (actual interface{}, loaded bool) {
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
func (m LinkedMap) Get(key interface{}) (value interface{}, loaded bool) {
	o, ok := m.m[key]
	if ok {
		return o.Value.([2]interface{})[1], ok
	}
	return nil, false
}

// 删除key对应的元素
// Param：key
func (m LinkedMap) Delete(key interface{}) {
	if e, ok := m.m[key]; ok {
		m.l.Remove(e)
		delete(m.m, key)
	}
}

// 获得Map长度
// Return： 链表长度
func (m LinkedMap) Size() int {
	return len(m.m)
}

// 轮询Map O(N)
// Param：接受轮询的函数，返回true继续轮询，返回false终止轮询
func (m LinkedMap) Foreach(f func(key interface{}, value interface{}) bool) {
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
func (m LinkedMap) Find(key interface{}) bool {
	_, ok := m.m[key]
	return ok
}
