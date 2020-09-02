/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description:
 */

package test

import (
	"github.com/xfali/goutils/net"
	"testing"
)

func TestNet(t *testing.T) {
	ret, _ := net.Format()
	for _, v := range ret {
		t.Logf("mac %s\n", v.Mac)
		for _, ip := range v.Ips {
			t.Logf("ip %s\n", ip)
		}
	}
}
