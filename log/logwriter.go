// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description:

package log

import (
	"errors"
	"os"
)

const (
	LogBufferSize = 10240
)

type FileLogWriter struct {
	stopChan chan bool
	logChan  chan []byte
	file     *os.File
}

func NewFileLogWriter(path string) *FileLogWriter {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil
	}
	l := FileLogWriter{
		stopChan: make(chan bool),
		logChan:  make(chan []byte, LogBufferSize),
		file:     f,
	}

	go func() {
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

func (w *FileLogWriter) writeLog(data []byte) {
	if w.file != nil {
		w.file.Write(data)
	}
}

func (w *FileLogWriter) Close() {
	close(w.stopChan)
	if w.file != nil {
		w.file.Close()
	}
}

func (w *FileLogWriter) Write(data []byte) (n int, err error) {
	select {
	case w.logChan <- data:
		return len(data), nil
	default:
		return 0, errors.New("write log failed ")
	}
}
