/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @date 2019/2/17
 * @time 10:15
 * @version V1.0
 * Description:
 */

package linkedSet

import "container/list"

type LinkedSet struct {
	l *list.List
	m map[interface{}] *list.Element
}

func New() *LinkedSet {
	return &LinkedSet{list.New(),make(map[interface{}] *list.Element)}
}

func (s *LinkedSet) PushBack(o interface{}, override bool) bool{
	if e, ok := s.m[o]; ok {
		if override {
			e.Value = o
		}
		return true
	}
	s.m[o] = s.l.PushBack(o)
	return false
}

func (s *LinkedSet) PushFront(o interface{}, override bool) bool{
	if e, ok := s.m[o]; ok {
		if override {
			e.Value = o
		}
		return true
	}
	s.m[o] = s.l.PushFront(o)
	return false
}

func (s *LinkedSet) Remove(o interface{}) {
	if e, ok := s.m[o]; ok {
		s.l.Remove(e)
		delete(s.m, e)
	}
}

func (s *LinkedSet) Front() interface{} {
	return s.l.Front().Value
}

func (s *LinkedSet) Back() interface{} {
	return s.l.Back().Value
}

func (s *LinkedSet) PopFront() interface{} {
	if s.l.Len() == 0 {
		return nil
	}
	e := s.l.Front()
	s.l.Remove(e)
	delete(s.m, e.Value)
	return e.Value
}

func (s *LinkedSet) PopBack() interface{} {
	if s.l.Len() == 0 {
		return nil
	}
	e := s.l.Back()
	s.l.Remove(e)
	delete(s.m, e.Value)
	return e.Value
}

func (s *LinkedSet) Len() int {
	return s.l.Len()
}

func (s *LinkedSet) Foreach(f func (interface{})bool) {
	for e:= s.l.Front(); e!= nil; e=e.Next() {
		if f(e.Value) {
			break
		}
	}
}

func (s *LinkedSet) Find(i interface{}) bool {
	if _, ok := s.m[i]; ok {
		return true
	} else {
		return false
	}
}

