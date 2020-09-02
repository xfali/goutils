// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description:

package log

import (
	"errors"
	"io"
	"sync"
)

const (
	BufferSize = 10240
)

type AsyncLogWriter struct {
	stopChan chan bool
	logChan  chan []byte
	w        io.Writer
	block    bool
	wait     sync.WaitGroup
	once     sync.Once
}

// 异步写的Writer，本身Write、Close方法线程安全，参数WriteCloser可以非线程安全
// Param： w - 实际写入的Writer, bufSize - 接收的最大长度, block - 如果为true，则当超出bufSize大小时Write方法阻塞，否则返回error
func NewAsyncWriter(w io.Writer, bufSize int, block bool) *AsyncLogWriter {
	if bufSize <= 0 {
		bufSize = BufferSize
	}
	l := AsyncLogWriter{
		stopChan: make(chan bool),
		logChan:  make(chan []byte, bufSize),
		w:        w,
		block:    block,
	}
	l.wait.Add(1)

	go func() {
		defer l.wait.Done()
		for {
			select {
			case <-l.stopChan:
				return
			case d, ok := <-l.logChan:
				if ok {
					l.writeLog(d)
				}
			}
		}
	}()
	return &l
}

func (w *AsyncLogWriter) writeLog(data []byte) {
	if w.w != nil {
		w.w.Write(data)
	}
}

func (w *AsyncLogWriter) Close() error {
	w.once.Do(func() {
		close(w.stopChan)
		w.wait.Wait()
	})
	return nil
}

func (w *AsyncLogWriter) Write(data []byte) (n int, err error) {
	if len(data) == 0 {
		return 0, nil
	}

	if w.block {
		w.logChan <- data
		return len(data), nil
	} else {
		select {
		case w.logChan <- data:
			return len(data), nil
		default:
			return 0, errors.New("write log failed ")
		}
	}
}
