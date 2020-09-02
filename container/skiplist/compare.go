// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package skiplist

func CompareInt(a, b interface{}) int {
	x1 := 0
	if a != nil {
		x1 = a.(int)
	}
	x2 := b.(int)
	return x1 - x2
}

func CompareInt64(a, b interface{}) int {
	var x1 int64 = 0
	if a != nil {
		x1 = a.(int64)
	}
	x2 := b.(int64)
	return int(x1 - x2)
}

func CompareFloat32(a, b interface{}) int {
	var x1 float32 = 0
	if a != nil {
		x1 = a.(float32)
	}
	x2 := b.(float32)
	if x1 > x2 {
		return 1
	} else if x2 > x1 {
		return -1
	} else {
		return 0
	}
}

func CompareFloat64(a, b interface{}) int {
	var x1 float64 = 0
	if a != nil {
		x1 = a.(float64)
	}
	x2 := b.(float64)
	if x1 > x2 {
		return 1
	} else if x2 > x1 {
		return -1
	} else {
		return 0
	}
}

func SetKeyCompareInt() Opt {
	return func(list *SkipList) {
		list.keyCmp = CompareInt
	}
}

func SetKeyCompareInt64() Opt {
	return func(list *SkipList) {
		list.keyCmp = CompareInt64
	}
}

func SetKeyCompareFloat32() Opt {
	return func(list *SkipList) {
		list.keyCmp = CompareFloat32
	}
}

func SetKeyCompareFloat64() Opt {
	return func(list *SkipList) {
		list.keyCmp = CompareFloat64
	}
}
