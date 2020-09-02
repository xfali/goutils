// +build windows

// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package io

import "syscall"

func SetInvisible(path string) error {
	namep, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return err
	}
	if err := syscall.SetFileAttributes(namep, syscall.FILE_ATTRIBUTE_HIDDEN); err != nil {
		return err
	}
	return nil
}
