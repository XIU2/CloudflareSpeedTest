package task

import (
	"bufio"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const defaultInputFile = "ip.txt"

var (
	// IPv6 IP version is 6
	IPv6 = false
	// TestAll test all ip
	TestAll = false
	// IPFile is the filename of IP Rangs
	IPFile = defaultInputFile
)

func randipEndWith(num int) uint8 {
	rand.Seed(time.Now().UnixNano())
	return uint8(rand.Intn(num))
}

func loadIPRanges() []*net.IPAddr {
	if IPFile == "" {
		IPFile = defaultInputFile
	}
	file, err := os.Open(IPFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	firstIPs := make([]*net.IPAddr, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ipString := scanner.Text()
		// 如果不含有 '/' 则代表不是 IP 段，而是一个单独的 IP，因此需要加上 /32 /128 子网掩码
		if !strings.Contains(ipString, "/") {
			mask := "/32"
			if IPv6 {
				mask = "/128"
			}
			ipString += mask
		}
		if IPv6 {
			//IPv6
			firstIPs = append(firstIPs, loadIPv6(ipString)...)
			continue
		}
		// IPv4
		firstIPs = append(firstIPs, loadIPv4(ipString)...)
	}
	return firstIPs
}

func loadIPv4(ipString string) (firstIPs []*net.IPAddr) {
	firstIP, IPRange, err := net.ParseCIDR(ipString)
	// fmt.Println(firstIP)
	// fmt.Println(IPRange)
	if err != nil {
		log.Fatal(err)
	}
	minIP, maxIP, hostNum := getCidrIPRange(ipString) // 获取 IP 最后一段最小值和最大值
	for IPRange.Contains(firstIP) {
		if TestAll { // 如果是测速全部 IP
			for i := minIP; i <= maxIP; i++ { // 遍历 IP 最后一段最小值到最大值
				firstIP[15] = i
				firstIPCopy := make([]byte, len(firstIP))
				copy(firstIPCopy, firstIP)
				firstIPs = append(firstIPs, &net.IPAddr{IP: firstIPCopy})
			}
		} else { // 随机 IP 的最后一段 0.0.0.X
			firstIP[15] = minIP + randipEndWith(hostNum)
			firstIPCopy := make([]byte, len(firstIP))
			copy(firstIPCopy, firstIP)
			firstIPs = append(firstIPs, &net.IPAddr{IP: firstIPCopy})
		}
		firstIP[14]++ // 0.0.(X+1).X
		if firstIP[14] == 0 {
			firstIP[13]++ // 0.(X+1).X.X
			if firstIP[13] == 0 {
				firstIP[12]++ // (X+1).X.X.X
			}
		}
	}
	return
}

func loadIPv6(ipString string) (firstIPs []*net.IPAddr) {
	firstIP, IPRange, err := net.ParseCIDR(ipString)
	// fmt.Println(firstIP)
	// fmt.Println(IPRange)
	if err != nil {
		log.Fatal(err)
	}
	var tempIP uint8
	for IPRange.Contains(firstIP) {
		//fmt.Println(firstIP)
		//fmt.Println(firstIP[0], firstIP[1], firstIP[2], firstIP[3], firstIP[4], firstIP[5], firstIP[6], firstIP[7], firstIP[8], firstIP[9], firstIP[10], firstIP[11], firstIP[12], firstIP[13], firstIP[14], firstIP[15])
		if !strings.Contains(ipString, "/128") {
			firstIP[15] = randipEndWith(255) // 随机 IP 的最后一段
			firstIP[14] = randipEndWith(255) // 随机 IP 的最后一段
		}
		firstIPCopy := make([]byte, len(firstIP))
		copy(firstIPCopy, firstIP)
		firstIPs = append(firstIPs, &net.IPAddr{IP: firstIPCopy})
		tempIP = firstIP[13]
		firstIP[13] += randipEndWith(255)
		if firstIP[13] >= tempIP {
			continue
		}
		tempIP = firstIP[12]
		firstIP[12] += randipEndWith(255)
		if firstIP[12] >= tempIP {
			continue
		}
		tempIP = firstIP[11]
		firstIP[11] += randipEndWith(255)
		if firstIP[11] >= tempIP {
			continue
		}
		tempIP = firstIP[10]
		firstIP[10] += randipEndWith(255)
		if firstIP[10] >= tempIP {
			continue
		}
		tempIP = firstIP[9]
		firstIP[9] += randipEndWith(255)
		if firstIP[9] >= tempIP {
			continue
		}
		tempIP = firstIP[8]
		firstIP[8] += randipEndWith(255)
		if firstIP[8] >= tempIP {
			continue
		}
		tempIP = firstIP[7]
		firstIP[7] += randipEndWith(255)
		if firstIP[7] >= tempIP {
			continue
		}
		tempIP = firstIP[6]
		firstIP[6] += randipEndWith(255)
		if firstIP[6] >= tempIP {
			continue
		}
		tempIP = firstIP[5]
		firstIP[5] += randipEndWith(255)
		if firstIP[5] >= tempIP {
			continue
		}
		tempIP = firstIP[4]
		firstIP[4] += randipEndWith(255)
		if firstIP[4] >= tempIP {
			continue
		}
		tempIP = firstIP[3]
		firstIP[3] += randipEndWith(255)
		if firstIP[3] >= tempIP {
			continue
		}
		tempIP = firstIP[2]
		firstIP[2] += randipEndWith(255)
		if firstIP[2] >= tempIP {
			continue
		}
		tempIP = firstIP[1]
		firstIP[1] += randipEndWith(255)
		if firstIP[1] >= tempIP {
			continue
		}
		tempIP = firstIP[0]
		firstIP[0] += randipEndWith(255)
	}
	return
}

// 根据子网掩码获取主机数量
func getCIDRHostNum(mask uint8) (subnetNum int) {
	if mask >= 32 {
		return 1
	}

	if mask < 32 {
		for i := int(32 - mask - 1); i >= 1; i-- {
			subnetNum += 1 << i
		}
		subnetNum += 2
	}
	if subnetNum > 0xFF {
		subnetNum = 0xFF
	}
	return
}

// 获取 IP 最后一段最小值和最大值、主机数量
func getCidrIPRange(cidr string) (minIP, maxIP uint8, ipNum int) {
	i := strings.IndexByte(cidr, '/')
	addr := cidr[:i]
	mask, _ := strconv.ParseUint(cidr[i+1:], 10, 8)
	i = strings.LastIndexByte(addr, '.')
	seg4, _ := strconv.ParseUint(addr[i+1:], 10, 8)
	minIP, maxIP = getIPSegRange(uint8(seg4), uint8(32-mask))
	ipNum = getCIDRHostNum(uint8(mask))
	return
}

// 根据输入的基础IP地址和CIDR掩码计算一个IP片段的区间
func getIPSegRange(userSegIP, offset uint8) (uint8, uint8) {
	var ipSegMax uint8 = 0xFF
	netSegIP := ipSegMax << offset
	segMinIP := netSegIP & userSegIP
	segMaxIP := userSegIP&(0xFF<<offset) | ^(0xFF << offset)
	return segMinIP, segMaxIP
}
