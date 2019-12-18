// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package log

import (
    "errors"
    "os"
    "time"
)

const (
    FlushSize     = 10240
    BufferNum     = 10
    FlushTime     = 5 * time.Millisecond
)

type FileBufferLogWriter struct {
    stopChan   chan bool
    logChan    chan []byte
    logBuffers [][]byte
    curBuf     int
    curSize    int
    flushTime  time.Time
    file       *os.File
}

func NewFileBufferLogWriter(path string) *FileBufferLogWriter {
    f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return nil
    }
    l := FileBufferLogWriter{
        stopChan:   make(chan bool),
        logChan:    make(chan []byte, LogBufferSize),
        logBuffers: make([][]byte, BufferNum),
        flushTime:  time.Now(),
        file:       f,
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

func (w *FileBufferLogWriter) writeLog(data []byte) {
    w.curSize += len(data)
    w.logBuffers[w.curBuf] = data
    w.curBuf++
    now := time.Now().Sub(w.flushTime)

    if w.curSize < FlushSize && w.curBuf < BufferNum && now < FlushTime {
        return
    }

    if w.file != nil {
        for i := range w.logBuffers {
            w.file.Write(w.logBuffers[i])
        }
    }
    w.curBuf = 0
    w.curSize = 0
    w.flushTime = time.Now()
}

func (w *FileBufferLogWriter) Close() {
    close(w.stopChan)
    if w.file != nil {
        w.file.Close()
    }
}

func (w *FileBufferLogWriter) Write(data []byte) (n int, err error) {
    if data == nil || len(data) == 0 {
        return 0, nil
    }
    select {
    case w.logChan <- data:
        return len(data), nil
    default:
        return 0, errors.New("write log failed ")
    }
}
