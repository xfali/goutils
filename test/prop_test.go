// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description: 

package test

import (
    "github.com/xfali/goutils/prop"
    "testing"
)

func TestEnv(t *testing.T) {
    ret := prop.GetEnvs()
    for k, v := range ret {
        t.Logf("key: %s , value: %s \n", k, v)
    }
}
