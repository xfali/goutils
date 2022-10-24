/*
 * Copyright 2022 Xiongfa Li.
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

package purger

import (
	"github.com/xfali/timewheel"
	"sync"
	"time"
)

// Purger 清理器
type Purger interface {
	// Purge 清理方法，需要保证线程安全，可重复调用
	Purge()
}

// PurgeExecutor 清理执行器
type PurgeExecutor interface {
	// AddPurger 添加清理器到清理执行器
	// purger为清理器，到interval时间间隔之后会自动调用清理器的Purge方法进行清理
	AddPurger(purger Purger, interval time.Duration) bool

	// Close 关闭清理执行器，关闭时会自动调用所有清理器再执行一次清理
	Close() error
}

type defaultExecutor struct {
	tw      timewheel.TimeWheel
	purgers map[Purger]struct{}
	lock    sync.Mutex
	once    sync.Once
}

type opt func(*defaultExecutor)

// New 创建默认的清理执行器
func New(opts ...opt) *defaultExecutor {
	ret := &defaultExecutor{
		purgers: make(map[Purger]struct{}),
		tw: timewheel.NewAsyncHiera(
			24*time.Hour,
			[]time.Duration{time.Hour, time.Minute, time.Second, 100 * time.Millisecond},
			100, 100),
	}
	for _, opt := range opts {
		opt(ret)
	}
	ret.tw.Start()
	return ret
}

func (e *defaultExecutor) AddPurger(purger Purger, interval time.Duration) bool {
	e.lock.Lock()
	defer e.lock.Unlock()

	if _, ok := e.purgers[purger]; ok {
		return false
	}
	_, err := e.tw.Add(purger.Purge, interval, true)
	if err != nil {
		return false
	}
	e.purgers[purger] = struct{}{}
	return true
}

func (e *defaultExecutor) Close() error {
	e.once.Do(func() {
		if e.tw != nil {
			e.tw.Stop()
		}
		e.lock.Lock()
		defer e.lock.Unlock()
		for p := range e.purgers {
			p.Purge()
		}
		e.purgers = nil
	})
	return nil
}

type options struct {
}

var Opts options

func (o options) TimeWheel(tw timewheel.TimeWheel) opt {
	return func(executor *defaultExecutor) {
		if tw != nil {
			executor.tw = tw
		}
	}
}
