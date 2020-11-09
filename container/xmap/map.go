// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package xmap

type IMap interface {
	// 向Map中添加一个元素
	// Param：key 添加的对象key，value 添加的对象
	Put(key, value interface{})

	// 获取key对应的元素
	// Param：key 对象key
	// Return： value：key对应的对象，loaded：成功获取返回true，不存在返回false
	Get(key interface{}) (value interface{}, loaded bool)

	// 删除key对应的元素
	// Param：key
	Delete(key interface{})

	// 获得Map长度
	// Return： 链表长度
	Size() int
}

type Map interface {
	IMap

	// 尝试向Map中添加一个元素，如果已存在该元素则直接返回已存在元素不进行添加
	// Param：key 添加的对象key，value 添加的对象
	// Return： actual 如果key已存在对应元素，则返回该元素，否则返回新添加的元素。 loaded：已存在返回true，否则返回false
	GetOrPut(key, value interface{}) (actual interface{}, loaded bool)

	// 轮询Map O(N)
	// Param：接受轮询的函数，返回true继续轮询，返回false终止轮询
	Foreach(f func(interface{}, interface{}) bool)
}
