// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package xlist

type List interface {
	// 在链表尾部添加一个元素
	// Param：o 添加的对象
	// Return： 成功返回true，失败返回false
	PushBack(o interface{}) bool

	// 在链表首部添加一个元素
	// Param：o 添加的对象
	// Return： 成功返回true，失败返回false
	PushFront(o interface{}) bool

	// 删除元素
	// Param：o 添加的对象
	Remove(o interface{})

	// 获得首元素
	// Return： 首元素（第一个元素），如果没有返回nil
	Front() interface{}

	// 获得尾元素
	// Return： 尾元素（最后一个元素），如果没有返回nil
	Back() interface{}

	// 获得首元素并移除
	// Return： 首元素（第一个元素），如果没有返回nil
	PopFront() interface{}

	// 获得尾元素并移除
	// Return： 尾元素（最后一个元素），如果没有返回nil
	PopBack() interface{}

	// 获得链表长度
	// Return： 链表长度
	Len() int

	// 轮询链表O(N)
	// Param：接受轮询的函数，返回true继续轮询，返回false终止轮询
	Foreach(f func(interface{}) bool)

	// 查询链表中是否存在参数对象
	// Param：查询的对象
	// Return：存在返回true，不存在返回false
	Find(i interface{}) bool
}
