// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package sortmap

import "errors"

type Iterator interface {
	HasNext() bool
	Next() (string, interface{})
}

type SortMap interface {
	Add(keyAndValues ...interface{}) error
	GetAll() map[string]interface{}
	Keys() []string
	Get(key string) interface{}
	Remove(key string) error
	Len() int

	Iterator() Iterator

	Clone() SortMap
}

type defaultSortMap [2]interface{}

func New(keyAndValues ...interface{}) SortMap {
	ret := &defaultSortMap{}
	ret[0] = []string{}
	ret[1] = map[string]interface{}{}

	ret.Add(keyAndValues...)
	return ret
}

func (f *defaultSortMap) Add(keyAndValues ...interface{}) error {
	size := len(keyAndValues)
	if size == 0 {
		return nil
	}
	//keys := f[0].([]string)
	kvs := f[1].(map[string]interface{})
	var k string
	for i := 0; i < size; i++ {
		if i%2 == 0 {
			if keyAndValues[i] == nil {
				return errors.New("Key must be not nil ")
			}
			s, ok := keyAndValues[i].(string)
			if !ok {
				return errors.New("Key must be string ")
			}
			k = s
			if _, ok := kvs[k]; !ok {
				f[0] = append(f[0].([]string), k)
			}
		} else {
			kvs[k] = keyAndValues[i]
		}
	}
	return nil
}

func (f *defaultSortMap) GetAll() map[string]interface{} {
	return f[1].(map[string]interface{})
}

func (f *defaultSortMap) Iterator() Iterator {
	return &defaultIterator{
		field: f,
		cur:   0,
	}
}

func (f defaultSortMap) Keys() []string {
	return f[0].([]string)
}

func (f defaultSortMap) Get(key string) interface{} {
	return f[1].(map[string]interface{})[key]
}

func (f *defaultSortMap) Remove(key string) error {
	_, ok := f[1].(map[string]interface{})[key]
	if ok {
		delete(f[1].(map[string]interface{}), key)
		keys := f[0].([]string)
		for i := 0; i < len(keys); i++ {
			if keys[i] == key {
				f[0] = append(keys[:i], keys[i+1:]...)
				break
			}
		}
		return nil
	} else {
		return errors.New("Key not found ")
	}
}

func (f defaultSortMap) Len() int {
	return len(f[0].([]string))
}

func (f defaultSortMap) Clone() SortMap {
	ret := New()

	for _, k := range f.Keys() {
		ret.Add(k, f.Get(k))
	}

	return ret
}

type defaultIterator struct {
	field *defaultSortMap
	cur   int
}

func (c *defaultIterator) HasNext() bool {
	return c.cur < len(c.field[0].([]string))
}

func (c *defaultIterator) Next() (string, interface{}) {
	v := c.field[0].([]string)[c.cur]
	c.cur++
	return v, c.field[1].(map[string]interface{})[v]
}

func MergeFields(fields ...SortMap) (SortMap, error) {
	if len(fields) == 0 {
		return nil, errors.New("No field to merge ")
	} else {
		field := fields[0]
		var tmp SortMap
		for i := 1; i < len(fields); i++ {
			tmp = fields[i]
			if tmp == nil {
				continue
			}
			keys := tmp.Keys()
			for _, k := range keys {
				err := field.Add(k, tmp.Get(k))
				if err != nil {
					return field, err
				}
			}
		}
		return field, nil
	}
}
