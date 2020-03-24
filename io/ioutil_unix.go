// +build !windows

// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description: 

package io

import "errors"

func SetInvisible(path string) error {
    return errors.New("Not support")
}
