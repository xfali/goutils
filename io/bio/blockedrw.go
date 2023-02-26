/*
 * Copyright (C) 2023, Xiongfa Li.
 * All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package bio

import (
	"bytes"
	"errors"
	"io"
)

type ReaderFactory func([]byte) io.Reader

type brwOpt func(*blockReadWriter)

type Pool interface {
	Put(x any)
	Get() any
}

type blockReadWriter struct {
	ch   chan io.Reader
	r    io.Reader
	rf   ReaderFactory
	pool Pool
}

func NewBlockReadWriter(opts ...brwOpt) *blockReadWriter {
	ret := &blockReadWriter{
		ch: make(chan io.Reader),
	}
	for _, opt := range opts {
		opt(ret)
	}
	return ret
}

func (rw *blockReadWriter) Read(d []byte) (int, error) {
	total := 0
	for {
		var rn int = 0
		var err error
		if rw.r != nil {
			for {
				rn, err = rw.r.Read(d)
				total += rn
				if errors.Is(err, io.EOF) {
					if rw.pool != nil {
						rw.pool.Put(rw.r)
					}
					break
				}
				if err != nil || rn == len(d) {
					return total, err
				}
				d = d[rn:]
			}
		}
		r, ok := <-rw.ch
		if !ok {
			return total, io.EOF
		}
		rw.r = r
	}
}

func (rw *blockReadWriter) Write(d []byte) (int, error) {
	if len(d) == 0 {
		return 0, nil
	}
	if rw.pool != nil {
		v := rw.pool.Get().(io.ReadWriter)
		n, err := v.Write(d)
		rw.ch <- v
		return n, err
	} else if rw.rf != nil {
		rw.ch <- rw.rf(d)
	} else {
		rw.ch <- bytes.NewReader(d)
	}
	return len(d), nil
}

func (rw *blockReadWriter) Close() error {
	select {
	case <-rw.ch:
		return nil
	default:
		close(rw.ch)
	}
	return nil
}

type brwOpts struct{}

var BlockRwOpts brwOpts

func (o brwOpts) WithReaderFactory(rf ReaderFactory) brwOpt {
	return func(writer *blockReadWriter) {
		writer.rf = rf
	}
}

func (o brwOpts) WithPool(p Pool) brwOpt {
	return func(writer *blockReadWriter) {
		writer.pool = p
	}
}
