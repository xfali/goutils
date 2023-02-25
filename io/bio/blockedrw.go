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
	"sync"
)

type ReaderFactory func([]byte) io.Reader

type brwOpt func(*blockedReadWriter)

type blockedReadWriter struct {
	ch   chan io.Reader
	r    io.Reader
	rf   ReaderFactory
	pool *sync.Pool
}

func NewBlockedReadWriter(opts ...brwOpt) *blockedReadWriter {
	ret := &blockedReadWriter{
		ch: make(chan io.Reader),
	}
	for _, opt := range opts {
		opt(ret)
	}
	return ret
}

func (rw *blockedReadWriter) Read(d []byte) (int, error) {
	total := 0
	for {
		var rn int = 0
		var err error
		if rw.r != nil {
			rn, err = rw.r.Read(d)
			total += rn
			if errors.Is(err, io.EOF) {
				if rw.pool != nil {
					rw.pool.Put(rw.r)
				}
			}

			if rn == len(d) {
				return total, nil
			}
		}
		r, ok := <-rw.ch
		if !ok {
			return total, io.EOF
		}
		rw.r = r
		d = d[rn:]
	}
}

func (rw *blockedReadWriter) Write(d []byte) (int, error) {
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

func (rw *blockedReadWriter) Close() error {
	select {
	case <-rw.ch:
		return nil
	default:
		close(rw.ch)
	}
	return nil
}
