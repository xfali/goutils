/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @date 2019/2/28
 * @time 17:16
 * @version V1.0
 * Description: 
 */

package id

import (
    "crypto/rand"
    "encoding/base64"
)

func RandomId(length int) string {
    b := make([]byte, length)
    if _, err := rand.Read(b); err != nil {
        return ""
    }
    return base64.URLEncoding.EncodeToString(b)
}
