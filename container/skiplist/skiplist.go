// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description: 

package skiplist

import (
	"math/rand"
)

const (
	SKIPLIST_MAX_LEVEL = 32
	SKIPLIST_P         = 0.25
)

type Compare func(a, b interface{}) int
type Opt func(*SkipList)

type Node struct {
	key     interface{}
	value   interface{}
	forward []*Node
}

type SkipList struct {
	maxLv int
	p     float32

	level  int
	len    int
	header *Node
	keyCmp Compare
}

func New(opts ...Opt) *SkipList {
	ret := &SkipList{
		maxLv: SKIPLIST_MAX_LEVEL,
		p:     SKIPLIST_P,
	}
	for _, v := range opts {
		v(ret)
	}

	if ret.keyCmp == nil {
		panic("keyCmp or valueCmp is nil!")
	}
	ret.header = makeNode(ret.maxLv, nil, nil)
	fillForward(ret.header.forward)
	return ret
}

func SetKeyCompareFunc(cmp Compare) Opt {
	return func(list *SkipList) {
		list.keyCmp = cmp
	}
}

func SetMaxLevel(maxLv int) Opt {
	return func(list *SkipList) {
		list.maxLv = maxLv
	}
}

func SetP(p float32) Opt {
	return func(list *SkipList) {
		list.p = p
	}
}

func randomLevel() int {
	level := 1
	f := 0xFFFF * float32(SKIPLIST_P)
	for ((rand.Int() & 0xFFFF) < int(f)) && (SKIPLIST_MAX_LEVEL > level) {
		level += 1
	}

	return level
}

func (list *SkipList) Search(searchKey interface{}) interface{} {
	x := list.header
	i := list.level
	for i >= 0 {
		for list.keyCmp(x.forward[i].key, searchKey) < 0 {
			x = x.forward[i]
		}
		i--
	}
	x = x.forward[0]
	if list.keyCmp(x.key, searchKey) == 0 {
		return x.value
	} else {
		return nil
	}
}

func (list *SkipList) Insert(searchKey, newValue interface{}) {
	update := make([]*Node, list.maxLv)
	x := list.header
	i := list.level
	for i >= 0 {
		for list.keyCmp(x.forward[i].key, searchKey) < 0 {
			x = x.forward[i]
		}
		update[i] = x
		i--
	}
	x = x.forward[0]
	if list.keyCmp(x.key, searchKey) == 0 {
		x.value = newValue
	} else {
		lvl := randomLevel()
		if lvl > list.level {
			i := list.level + 1
			for i <= lvl {
				update[i] = list.header
				i++
			}
			list.level = lvl
		}
		x = makeNode(lvl, searchKey, newValue)
		i := 0
		for i <= lvl {
			x.forward[i] = update[i].forward[i]
			update[i].forward[i] = x
			i++
		}
		list.len++
	}
}

func (list *SkipList) Delete(searchKey interface{}) bool {
	update := make([]*Node, list.maxLv)
	x := list.header
	i := list.level
	for i >= 0 {
		for list.keyCmp(x.forward[i].key, searchKey) < 0 {
			x = x.forward[i]
		}
		update[i] = x
		i--
	}
	x = x.forward[0]
	if list.keyCmp(x.key, searchKey) == 0 {
		i := 0
		for i <= list.level {
			if update[i].forward[i] != x {
				break
			}
			update[i].forward[i] = x.forward[i]
			i++
		}
		freeNode(x)
		for list.level > 0 && list.header.forward[list.level] == nil {
			list.level--
		}
		list.len--
		return true
	}
	return false
}

func (list *SkipList) Len() int {
	return list.len
}

func (list *SkipList) Values(size int) []interface{} {
	if size < 0 {
		size = list.len
	}

	if size == 0 {
		return nil
	}

	ret := make([]interface{}, size)
	i := 0
	x := list.header.forward[0]
	for x != nil {
		if i >= size {
			break
		}
		ret[i] = x.value
		x = x.forward[0]
		i++
	}
	return ret
}

func makeNode(lvl int, searchKey, newValue interface{}) *Node {
	return &Node{
		key:     searchKey,
		value:   newValue,
		forward: make([]*Node, lvl+1),
	}
}

func fillForward(nodes []*Node) {
	vs := make([]Node, len(nodes))
	for i := range nodes {
		nodes[i] = &vs[i]
	}
}

func freeNode(n *Node) {

}
