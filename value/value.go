// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package value

import (
	"sync"
	"sync/atomic"
)

// 存储值对象工具，interface不做类型检查，需用户自行确认存取类型
type Value interface {
	// 存储值
	Store(interface{})
	// 取出值
	Load() interface{}
}

type SimpleValue struct {
	o interface{}
}

func NewSimpleValue(o interface{}) *SimpleValue {
	return &SimpleValue{o: o}
}

func (l *SimpleValue) Load() interface{} {
	return l.o
}

func (l *SimpleValue) Store(o interface{}) {
	l.o = o
}

type AtomicValue atomic.Value

func NewAtomicValue(o interface{}) *AtomicValue {
	ret := &AtomicValue{}
	ret.Store(o)
	return ret
}
func (l *AtomicValue) Load() interface{} {
	return (*atomic.Value)(l).Load().(*SimpleValue).Load()
}

func (l *AtomicValue) Store(o interface{}) {
	(*atomic.Value)(l).Store(&SimpleValue{o: o})
}

type LockedValue struct {
	o    interface{}
	lock sync.RWMutex
}

func NewLockedValue(o interface{}) *LockedValue {
	ret := &LockedValue{o: o}
	return ret
}
func (l *LockedValue) Load() interface{} {
	l.lock.RLock()
	defer l.lock.RUnlock()

	return l.o
}

func (l *LockedValue) Store(o interface{}) {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.o = o
}
