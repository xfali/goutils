// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package lru

import "container/list"

type QueueElem list.Element
type QueueListener interface {
	// 元素命中，将元素移动到队列头部后调用
	PostTouch(v interface{})

	// 元素新添加到队列头部后调用
	PostInsert(v interface{})

	// 元素从队列中删除后调用
	PostDelete(v interface{})
}

type Queue struct {
	list     *list.List
	cap      int
	listener QueueListener
}

func NewLruElement(v interface{}) *QueueElem {
	return &QueueElem{
		Value: v,
	}
}

func NewLruQueue(cap int, listener QueueListener) *Queue {
	if listener == nil {
		listener = &dummyListener{}
	}
	return &Queue{
		list:     list.New(),
		cap:      cap,
		listener: listener,
	}
}

func (q *Queue) Touch(elem *QueueElem) {
	if elem != nil {
		q.list.MoveToFront((*list.Element)(elem))
		q.listener.PostTouch(elem.Value)
	}
}

func (q *Queue) Insert(v interface{}) *QueueElem {
	if q.list.Len() == q.cap {
		e := q.list.Back()
		if e != nil {
			q.Delete((*QueueElem)(e))
		}
	}

	e := q.list.PushFront(v)
	q.listener.PostInsert(v)
	return (*QueueElem)(e)
}

func (q *Queue) Delete(elem *QueueElem) {
	q.list.Remove((*list.Element)(elem))
	q.listener.PostDelete(elem.Value)
}

type dummyListener struct{}

func (l *dummyListener) PostTouch(v interface{}) {
}

func (l *dummyListener) PostInsert(v interface{}) {
}

func (l *dummyListener) PostDelete(v interface{}) {
}
