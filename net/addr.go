/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @date 2019/2/28
 * @time 17:45
 * @version V1.0
 * Description: 
 */

package net

import "net"

type Info struct {
    Mac string
    Ips []string
}

func Format() ([]Info, error) {
    intes, err := net.Interfaces()
    if err == nil {
        var ret []Info
        for i:=0; i<len(intes); i++ {
            mac := intes[i].HardwareAddr.String()
            if mac == "" {
                continue
            }

            tmp := Info{
                Mac: mac,
            }

            if addrs, err := intes[i].Addrs(); err == nil {
                for j:=0; j<len(addrs); j++ {
                    tmp.Ips = append(tmp.Ips, addrs[j].String())
                }
            }
            ret = append(ret, tmp)
        }
        return ret, nil
    }
    return nil, err
}
