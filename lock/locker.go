// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package lock

type FakeLocker struct{}

func (FakeLocker) Lock()   {}
func (FakeLocker) Unlock() {}
