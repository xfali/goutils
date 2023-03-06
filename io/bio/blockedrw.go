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
	"context"
	"errors"
	"io"
	"time"
)

type ReaderFactory func([]byte) io.Reader

type brwOpt func(*blockReadWriter)

type Pool interface {
	Put(x any)
	Get() any
}

type blockReadWriter struct {
	ch      chan io.Reader
	closeCh chan struct{}
	r       io.Reader
	rf      ReaderFactory
	pool    Pool

	readTimeout  time.Duration
	writeTimeout time.Duration
}

func NewBlockReadWriter(opts ...brwOpt) *blockReadWriter {
	ret := &blockReadWriter{
		ch:      make(chan io.Reader),
		closeCh: make(chan struct{}),
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
		var (
			r  io.Reader
			ok bool
		)
		if rw.readTimeout > 0 {
			ctx, _ := context.WithTimeout(context.Background(), rw.readTimeout)
			select {
			case <-rw.closeCh:
				select {
				case r, ok = <-rw.ch:
					break
				default:
					return total, io.EOF
				}
			case <-ctx.Done():
				return total, errors.New("Read timeout ")
			case r, ok = <-rw.ch:
				break
			}
		} else {
			select {
			case <-rw.closeCh:
				select {
				case r, ok = <-rw.ch:
					break
				default:
					return total, io.EOF
				}
			case r, ok = <-rw.ch:
				break
			}
		}

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
	var (
		r   io.Reader
		n   int
		err error
	)
	if rw.pool != nil {
		v := rw.pool.Get().(io.ReadWriter)
		n, err = v.Write(d)
		if err != nil {
			return n, err
		}
		r = v
	} else if rw.rf != nil {
		r = rw.rf(d)
		n = len(d)
	} else {
		r = bytes.NewReader(d)
		n = len(d)
	}

	if rw.writeTimeout > 0 {
		ctx, _ := context.WithTimeout(context.Background(), rw.writeTimeout)
		select {
		case <-rw.closeCh:
			return 0, errors.New("Closed ")
		case <-ctx.Done():
			return 0, errors.New("Write timeout ")
		case rw.ch <- r:
			break
		}
	} else {
		select {
		case <-rw.closeCh:
			return 0, errors.New("Closed ")
		case rw.ch <- r:
			break
		}
	}
	return n, err
}

func (rw *blockReadWriter) Close() error {
	select {
	case <-rw.closeCh:
		return nil
	default:
		close(rw.closeCh)
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

func (o brwOpts) WithReadTimeout(t time.Duration) brwOpt {
	return func(writer *blockReadWriter) {
		writer.readTimeout = t
	}
}

func (o brwOpts) WithWriteTimeout(t time.Duration) brwOpt {
	return func(writer *blockReadWriter) {
		writer.writeTimeout = t
	}
}
