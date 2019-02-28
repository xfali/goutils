/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @date 2019/2/28
 * @time 17:20
 * @version V1.0
 * Description: 
 */

package test

import (
    "goutils/id"
    "testing"
)

func TestRandomId(t *testing.T) {
    t.Logf("%s\n", id.RandomId(10))
    t.Logf("%s\n", id.RandomId(32))
    t.Logf("%s\n", id.RandomId(64))
}
