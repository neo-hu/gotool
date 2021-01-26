package net

import (
	"net"
	"regexp"
	"strconv"
	"strings"
)

func LocalIP() (ips []net.IP, err error) {
	ips = make([]net.IP, 0)
	ifaces, e := net.Interfaces()
	if e != nil {
		return ips, e
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}

		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}

		if strings.HasPrefix(iface.Name, "docker") || strings.HasPrefix(iface.Name, "w-") {
			continue
		}

		addrs, e := iface.Addrs()
		if e != nil {
			return ips, e
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}
			ips = append(ips, ip)
		}
	}
	return ips, nil
}

func IsIntranet(ipStr string) bool {
	if strings.HasPrefix(ipStr, "10.") || strings.HasPrefix(ipStr, "192.168.") {
		return true
	}

	if strings.HasPrefix(ipStr, "172.") {
		// 172.16.0.0-172.31.255.255
		arr := strings.Split(ipStr, ".")
		if len(arr) != 4 {
			return false
		}

		second, err := strconv.ParseInt(arr[1], 10, 64)
		if err != nil {
			return false
		}

		if second >= 16 && second <= 31 {
			return true
		}
	}

	return false
}

var ipRe *regexp.Regexp

func init() {
	ipRe = regexp.MustCompile(`^(?:(?:2[0-4][0-9]\.)|(?:25[0-5]\.)|(?:1[0-9][0-9]\.)|(?:[1-9][0-9]\.)|(?:[0-9]\.)){3}(?:(?:2[0-5][0-5])|(?:25[0-5])|(?:1[0-9][0-9])|(?:[1-9][0-9])|(?:[0-9]))$`)
}

func IpVerification(ip string) bool {
	return ipRe.MatchString(ip)
}
