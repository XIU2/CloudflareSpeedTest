package task

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/XIU2/CloudflareSpeedTest/utils"
)

var (
	NetworkInterface string // 指定网卡名称或本地源 IP
	boundInterface   *net.Interface
	boundIPv4        net.IP
	boundIPv6        net.IP
)

// InitInterface 初始化测速指定的网卡或本地 IP
func InitInterface() {
	if NetworkInterface == "" {
		return
	}

	// 1. 如果用户传入的是 IP 地址
	if ip := net.ParseIP(NetworkInterface); ip != nil {
		if ip.To4() != nil {
			boundIPv4 = ip
		} else {
			boundIPv6 = ip
		}
		// 尝试匹配该 IP 所属的网络接口
		if ifaces, err := net.Interfaces(); err == nil {
			for i := range ifaces {
				addrs, err := ifaces[i].Addrs()
				if err != nil {
					continue
				}
				for _, addr := range addrs {
					var ifaceIP net.IP
					switch v := addr.(type) {
					case *net.IPNet:
						ifaceIP = v.IP
					case *net.IPAddr:
						ifaceIP = v.IP
					}
					if ifaceIP != nil && ifaceIP.Equal(ip) {
						ifaceCopy := ifaces[i]
						boundInterface = &ifaceCopy
						break
					}
				}
				if boundInterface != nil {
					break
				}
			}
		}
		if boundInterface != nil {
			utils.Cyan.Printf("已指定测速网卡: %s (本地 IP: %s)\n", boundInterface.Name, ip.String())
		} else {
			utils.Cyan.Printf("已指定测速本地 IP: %s\n", ip.String())
		}
		return
	}

	// 2. 如果用户传入的是网卡名称 (如 en1, eth0)
	iface, err := net.InterfaceByName(NetworkInterface)
	if err != nil {
		var avail []string
		if ifaces, err2 := net.Interfaces(); err2 == nil {
			for _, i := range ifaces {
				avail = append(avail, i.Name)
			}
		}
		log.Fatalf("未找到指定的网卡或本地 IP: [%s]，当前可用网卡: %s\n", NetworkInterface, strings.Join(avail, ", "))
	}
	boundInterface = iface

	// 获取该网卡上的 IPv4 和 IPv6 地址
	if addrs, err := iface.Addrs(); err == nil {
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil {
				continue
			}
			if ip.IsLoopback() && (iface.Flags&net.FlagLoopback == 0) {
				continue
			}
			if ip.To4() != nil {
				if boundIPv4 == nil {
					boundIPv4 = ip
				}
			} else {
				if boundIPv6 == nil && !ip.IsLinkLocalUnicast() {
					boundIPv6 = ip
				}
			}
		}
	}

	out := fmt.Sprintf("已指定测速网卡: %s", boundInterface.Name)
	if boundIPv4 != nil {
		out += fmt.Sprintf(" (IPv4: %s)", boundIPv4.String())
	}
	if boundIPv6 != nil {
		out += fmt.Sprintf(" (IPv6: %s)", boundIPv6.String())
	}
	utils.Cyan.Println(out)
}

// newDialer 创建配置了超时和指定网卡绑定的 net.Dialer
func newDialer(ip *net.IPAddr, timeout time.Duration) *net.Dialer {
	dialer := &net.Dialer{
		Timeout: timeout,
	}
	if NetworkInterface == "" {
		return dialer
	}

	targetIsIPv4 := true
	if ip != nil {
		targetIsIPv4 = isIPv4(ip.String())
	}

	if targetIsIPv4 && boundIPv4 != nil {
		dialer.LocalAddr = &net.TCPAddr{IP: boundIPv4}
	} else if !targetIsIPv4 && boundIPv6 != nil {
		dialer.LocalAddr = &net.TCPAddr{IP: boundIPv6}
	}

	if boundInterface != nil {
		dialer.Control = bindInterfaceControl(boundInterface, targetIsIPv4)
	}

	return dialer
}
