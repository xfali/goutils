// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package stack

import "container/list"

type Stack struct {
	l *list.List
}

func New() *Stack {
	return &Stack{
		l: list.New(),
	}
}

func (s *Stack) Push(v interface{}) {
	s.l.PushBack(v)
}

func (s *Stack) Pop() interface{} {
	e := s.l.Back()
	if e == nil {
		return nil
	}

	return s.l.Remove(e)
}

func (s *Stack) Len() int {
	return s.l.Len()
}

func (s *Stack) Foreach(f func(interface{}) bool) {
	for e := s.l.Front(); e != nil; e = e.Next() {
		if f(e.Value) {
			break
		}
	}
}
