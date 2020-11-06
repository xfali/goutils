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


