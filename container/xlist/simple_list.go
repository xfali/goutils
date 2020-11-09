// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package xlist

import (
	"container/list"
	"reflect"
)

type SimpleList struct {
	l *list.List
}

func NewSimpleList() *SimpleList {
	return &SimpleList{
		l: list.New(),
	}
}

// 在链表尾部添加一个元素
// Param：o 添加的对象
func (l *SimpleList) PushBack(o interface{}) {
	l.l.PushBack(o)
}

// 在链表首部添加一个元素
// Param：o 添加的对象
func (l *SimpleList) PushFront(o interface{}) {
	l.l.PushFront(o)
}

// 删除元素O(N)
// Param：o 添加的对象
func (l *SimpleList) Remove(o interface{}) {
	for e := l.l.Front(); e != nil; e = e.Next() {
		if reflect.DeepEqual(o, e.Value) {
			l.l.Remove(e)
			return
		}
	}
}

// 获得首元素
// Return： 首元素（第一个元素），如果没有返回nil
func (l *SimpleList) Front() interface{} {
	e := l.l.Front()
	if e != nil {
		return e.Value
	}
	return nil
}

// 获得尾元素
// Return： 尾元素（最后一个元素），如果没有返回nil
func (l *SimpleList) Back() interface{} {
	e := l.l.Back()
	if e != nil {
		return e.Value
	}
	return nil
}

// 获得首元素并移除
// Return： 首元素（第一个元素），如果没有返回nil
func (l *SimpleList) PopFront() interface{} {
	e := l.l.Front()
	if e != nil {
		l.l.Remove(e)
		return e.Value
	}
	return nil
}

// 获得尾元素并移除
// Return： 尾元素（最后一个元素），如果没有返回nil
func (l *SimpleList) PopBack() interface{} {
	e := l.l.Back()
	if e != nil {
		l.l.Remove(e)
		return e.Value
	}
	return nil
}

// 获得链表长度
// Return： 链表长度
func (l *SimpleList) Len() int {
	return l.l.Len()
}

// 轮询链表O(N)
// Param：接受轮询的函数，返回true继续轮询，返回false终止轮询
func (l *SimpleList) Foreach(f func(interface{}) bool) {
	for e := l.l.Front(); e != nil; e = e.Next() {
		if !f(e.Value) {
			return
		}
	}
}

// 查询链表中是否存在参数对象O(N)
// Param：查询的对象
// Return：存在返回true，不存在返回false
func (l *SimpleList) Find(i interface{}) bool {
	for e := l.l.Front(); e != nil; e = e.Next() {
		if reflect.DeepEqual(i, e.Value) {
			return true
		}
	}
	return false
}
