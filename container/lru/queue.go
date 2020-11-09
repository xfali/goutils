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

type Queue interface {
	AddListener(listener QueueListener)

	Touch(elem *QueueElem)

	Move(other *LruQueue, elem *QueueElem, notify bool) *QueueElem

	Insert(v interface{}) *QueueElem

	Delete(elem *QueueElem)
}

type LruQueue struct {
	list      *list.List
	cap       int
	listeners []QueueListener
}

func NewLruElement(v interface{}) *QueueElem {
	return &QueueElem{
		Value: v,
	}
}

func NewLruQueue(cap int) *LruQueue {
	return &LruQueue{
		list:      list.New(),
		cap:       cap,
		listeners: []QueueListener{},
	}
}

func (q *LruQueue) AddListener(listener QueueListener) {
	q.listeners = append(q.listeners, listener)
}

func (q *LruQueue) Touch(elem *QueueElem) {
	if elem != nil {
		q.list.MoveToFront((*list.Element)(elem))
		for _, l := range q.listeners {
			l.PostTouch(elem.Value)
		}
	}
}

func (q *LruQueue) Move(other *LruQueue, elem *QueueElem, notify bool) *QueueElem {
	if other == nil || elem == nil {
		return nil
	}
	v := q.list.Remove((*list.Element)(elem))
	if notify {
		for _, l := range q.listeners {
			l.PostDelete(elem.Value)
		}
	}

	elem = (*QueueElem)(other.list.PushFront(v))
	if notify {
		for _, l := range q.listeners {
			l.PostInsert(elem.Value)
		}
	}
	return elem
}

func (q *LruQueue) Insert(v interface{}) *QueueElem {
	if q.list.Len() == q.cap {
		e := q.list.Back()
		if e != nil {
			q.Delete((*QueueElem)(e))
		}
	}

	e := q.list.PushFront(v)
	for _, l := range q.listeners {
		l.PostInsert(v)
	}
	return (*QueueElem)(e)
}

func (q *LruQueue) Delete(elem *QueueElem) {
	v := q.list.Remove((*list.Element)(elem))
	for _, l := range q.listeners {
		l.PostDelete(v)
	}
}

type dummyListener struct{}

func (l *dummyListener) PostTouch(v interface{}) {
}

func (l *dummyListener) PostInsert(v interface{}) {
}

func (l *dummyListener) PostDelete(v interface{}) {
}
