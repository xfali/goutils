// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package xmap

type SimpleMap map[interface{}]interface{}

func NewSimpleMap() *SimpleMap {
	return &SimpleMap{}
}

// 向Map中添加一个元素
// Param：key 添加的对象key，value 添加的对象
// Return： 成功返回true，失败返回false
func (m SimpleMap) Put(key, value interface{}) {
	m[key] = value
}

// 尝试向Map中添加一个元素，如果已存在该元素则直接返回已存在元素不进行添加
// Param：key 添加的对象key，value 添加的对象
// Return： actual 如果key已存在对应元素，则返回该元素，否则返回新添加的元素。 loaded：已存在返回true，否则返回false
func (m SimpleMap) GetOrPut(key, value interface{}) (actual interface{}, loaded bool) {
	o, ok := m[key]
	if ok {
		return o, true
	}
	m[key] = value
	return value, false
}

// 获取key对应的元素
// Param：key 对象key
// Return： value：key对应的对象，loaded：成功获取返回true，不存在返回false
func (m SimpleMap) Get(key interface{}) (value interface{}, loaded bool) {
	o, ok := m[key]
	if ok {
		return o, true
	}
	return nil, false
}

// 删除key对应的元素
// Param：key
func (m SimpleMap) Delete(key interface{}) {
	delete(m, key)
}

// 获得Map长度
// Return： 链表长度
func (m SimpleMap) Size() int {
	return len(m)
}

// 轮询Map O(N)
// Param：接受轮询的函数，返回true继续轮询，返回false终止轮询
func (m SimpleMap) Foreach(f func(interface{}, interface{}) bool) {
	for k, v := range m {
		if !f(k, v) {
			break
		}
	}
}

// 查询Map中是否存在参数对象
// Param：查询的对象
// Return：存在返回true，不存在返回false
func (m SimpleMap) Find(key interface{}) bool {
	_, ok := m[key]
	return ok
}
