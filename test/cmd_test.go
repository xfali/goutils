// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description: 

package test

import (
    "github.com/xfali/goutils/cmd"
    "testing"
)

func TestCmd(t *testing.T) {
    err := cmd.ExecCommand("echo", "test")
    if err != nil {
        t.Fatal(err)
    }
}
