/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package test

import (
    "fmt"
    "goutils/idUtil"
    "testing"
)

func TestRandomId(t *testing.T) {
    t.Logf("%s\n", idUtil.RandomId(10))
    t.Logf("%s\n", idUtil.RandomId(32))
    t.Logf("%s\n", idUtil.RandomId(64))
}

func TestSnowFlakeId(t *testing.T) {
    sf := idUtil.NewSnowFlake()
    id, err := sf.NextId()
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    fmt.Printf("idUtil is %d, timestamp %s limitstr %s\n", id.Int64(), id.Timestamp(), id.LimitString(30))
    fmt.Printf("time is %v\n", id.Time())

    str := idUtil.Compress2StringUL2(id.Int64(), 20)
    fmt.Printf("func compress str is %s\n", str)
    id2 := idUtil.Uncompress2LongUL(str)
    fmt.Printf("func uncompress %d\n", id2)

    sid := id.Compress()
    fmt.Printf("compress id %s\n", sid)
    fmt.Printf("uncompress id %d\n", sid.UnCompress())
}
