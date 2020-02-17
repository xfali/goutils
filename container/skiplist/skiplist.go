// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description: 

package skiplist

import (
	"math/rand"
	"time"
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

type Iterator struct {
	n *Node
}

type SkipList struct {
	maxLv int
	p     float32

	level int
	len   int

	rand   *rand.Rand
	header *Node
	keyCmp Compare
}

func New(opts ...Opt) *SkipList {
	ret := &SkipList{
		maxLv: SKIPLIST_MAX_LEVEL,
		p:     SKIPLIST_P,
		rand:  rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	for _, v := range opts {
		v(ret)
	}

	if ret.keyCmp == nil {
		panic("keyCmp or valueCmp is nil!")
	}
	ret.header = makeNode(ret.maxLv, nil, nil)
	//fillForward(ret.header.forward)
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

func randomLevel(rand *rand.Rand) int {
	level := 1
	f := 0xFFFF * float32(SKIPLIST_P)
	for ((rand.Int() & 0xFFFF) < int(f)) && (SKIPLIST_MAX_LEVEL > level) {
		level += 1
	}

	return level
}

func (list *SkipList) Get(searchKey interface{}) interface{} {
	x := list.header
	i := list.level
	for i >= 0 {
		for x.forward[i] != nil && list.keyCmp(x.forward[i].key, searchKey) < 0 {
			x = x.forward[i]
		}
		i--
	}
	x = x.forward[0]
	if x != nil && list.keyCmp(x.key, searchKey) == 0 {
		return x.value
	} else {
		return nil
	}
}

func (list *SkipList) Set(searchKey, newValue interface{}) {
	update := make([]*Node, list.maxLv)
	x := list.header
	i := list.level
	for i >= 0 {
		for x.forward[i] != nil && list.keyCmp(x.forward[i].key, searchKey) < 0 {
			x = x.forward[i]
		}
		update[i] = x
		i--
	}
	x = x.forward[0]
	if x != nil && list.keyCmp(x.key, searchKey) == 0 {
		x.value = newValue
	} else {
		lvl := randomLevel(list.rand)
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
		for x.forward[i] != nil && list.keyCmp(x.forward[i].key, searchKey) < 0 {
			x = x.forward[i]
		}
		update[i] = x
		i--
	}
	x = x.forward[0]
	if x != nil && list.keyCmp(x.key, searchKey) == 0 {
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

func (list *SkipList) First() *Iterator {
	if list.header.forward[0] == nil {
		return nil
	}
	return &Iterator{n: list.header.forward[0]}
}

func (list *SkipList) Last() *Iterator {
	x := list.header.forward[0]
	if x == nil {
		return nil
	}
	for x != nil {
		if x.forward[0] == nil {
			return &Iterator{n: x}
		}
		x = x.forward[0]
	}
	return nil
}

func (list *SkipList) FirstNear(searchKey interface{}) *Iterator {
	x := list.header
	i := list.level
	for i >= 0 {
		for x.forward[i] != nil && list.keyCmp(x.forward[i].key, searchKey) < 0 {
			x = x.forward[i]
		}
		i--
	}
	x = x.forward[0]
	if x == nil {
		return nil
	}
	return &Iterator{n: x}
}

func (it *Iterator) Next() *Iterator {
	if it.n == nil {
		return nil
	}
	it.n = it.n.forward[0]
	if it.n == nil {
		return nil
	}
	return it
}

func (it *Iterator) Key() interface{} {
	if it.n == nil {
		return nil
	}

	return it.n.key
}

func (it *Iterator) Value() interface{} {
	if it.n == nil {
		return nil
	}

	return it.n.value
}

func (it *Iterator) KeyValue() (interface{}, interface{}) {
	if it.n == nil {
		return nil, nil
	}

	return it.n.value, it.n.value
}

func makeNode(lvl int, searchKey, newValue interface{}) *Node {
	return &Node{
		key:     searchKey,
		value:   newValue,
		forward: make([]*Node, lvl+1),
	}
}

func freeNode(n *Node) {

}
