// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package test

import (
    "github.com/xfali/goutils/log"
    "testing"
    "time"
)

func TestLog(t *testing.T) {
    log.Info("123")
    log.Info("123")
    //log.Fatal("exit")
}

func TestLog2(t *testing.T) {
    log.Info("", 5, 1.1, "测试", time.Now())
    log.Info("", 5, 1.1, "测试", time.Now())
}
