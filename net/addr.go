/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
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
		for i := 0; i < len(intes); i++ {
			mac := intes[i].HardwareAddr.String()
			if mac == "" {
				continue
			}

			tmp := Info{
				Mac: mac,
			}

			if addrs, err := intes[i].Addrs(); err == nil {
				for j := 0; j < len(addrs); j++ {
					tmp.Ips = append(tmp.Ips, addrs[j].String())
				}
			}
			ret = append(ret, tmp)
		}
		return ret, nil
	}
	return nil, err
}

func IP2Long(ip net.IP) uint32 {
	a := uint32(ip[12])
	b := uint32(ip[13])
	c := uint32(ip[14])
	d := uint32(ip[15])
	return uint32(a<<24 | b<<16 | c<<8 | d)
}

func Long2IP(ip uint32) net.IP {
	a := byte((ip >> 24) & 0xFF)
	b := byte((ip >> 16) & 0xFF)
	c := byte((ip >> 8) & 0xFF)
	d := byte(ip & 0xFF)
	return net.IPv4(a, b, c, d)
}
