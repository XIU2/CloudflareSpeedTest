package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

// 根据子网掩码获取主机数量
func getCidrHostNum(maskLen int) int {
	cidrIPNum := int(0)
	if maskLen < 32 {
		var i int = int(32 - maskLen - 1)
		for ; i >= 1; i-- {
			cidrIPNum += 1 << i
		}
		cidrIPNum += 2
	} else {
		cidrIPNum = 1
	}
	if cidrIPNum > 255 {
		cidrIPNum = 255
	}
	return cidrIPNum
}

// 获取 IP 最后一段最小值和最大值
func getCidrIPRange(cidr string) (uint8, uint8) {
	ip := strings.Split(cidr, "/")[0]
	ipSegs := strings.Split(ip, ".")
	maskLen, _ := strconv.Atoi(strings.Split(cidr, "/")[1])
	seg4MinIP, seg4MaxIP := getIPSeg4Range(ipSegs, maskLen)
	//ipPrefix := ipSegs[0] + "." + ipSegs[1] + "." + ipSegs[2] + "."

	return seg4MinIP,
		seg4MaxIP
}

// 获取 IP 最后一段的区间
func getIPSeg4Range(ipSegs []string, maskLen int) (uint8, uint8) {
	ipSeg, _ := strconv.Atoi(ipSegs[3])
	return getIPSegRange(uint8(ipSeg), uint8(32-maskLen))
}

// 根据输入的基础IP地址和CIDR掩码计算一个IP片段的区间
func getIPSegRange(userSegIP, offset uint8) (uint8, uint8) {
	var ipSegMax uint8 = 255
	netSegIP := ipSegMax << offset
	segMinIP := netSegIP & userSegIP
	segMaxIP := userSegIP&(255<<offset) | ^(255 << offset)
	return uint8(segMinIP), uint8(segMaxIP)
}

func loadFirstIPOfRangeFromFile(ipFile string) []net.IPAddr {
	file, err := os.Open(ipFile)
	if err != nil {
		log.Fatal(err)
	}
	firstIPs := make([]net.IPAddr, 0)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		IPString := scanner.Text()
		if !strings.Contains(IPString, "/") { // 如果不含有 / 则代表不是 IP 段，而是一个单独的 IP，因此需要加上 /32 子网掩码
			IPString += "/32"
		}
		firstIP, IPRange, err := net.ParseCIDR(IPString)
		//fmt.Println(firstIP)
		//fmt.Println(IPRange)
		if err != nil {
			log.Fatal(err)
		}
		if !ipv6Mode { // IPv4
			minIP, maxIP := getCidrIPRange(IPString)                 // 获取 IP 最后一段最小值和最大值
			Mask, _ := strconv.Atoi(strings.Split(IPString, "/")[1]) // 获取子网掩码
			MaxIPNum := getCidrHostNum(Mask)                         // 根据子网掩码获取主机数量
			for IPRange.Contains(firstIP) {
				if allip { // 如果是测速全部 IP
					for i := int(minIP); i <= int(maxIP); i++ { // 遍历 IP 最后一段最小值到最大值
						firstIP[15] = uint8(i)
						firstIPCopy := make([]byte, len(firstIP))
						copy(firstIPCopy, firstIP)
						firstIPs = append(firstIPs, net.IPAddr{IP: firstIPCopy})
					}
				} else { // 随机 IP 的最后一段 0.0.0.X
					firstIP[15] = minIP + randipEndWith(MaxIPNum)
					firstIPCopy := make([]byte, len(firstIP))
					copy(firstIPCopy, firstIP)
					firstIPs = append(firstIPs, net.IPAddr{IP: firstIPCopy})
				}
				firstIP[14]++ // 0.0.(X+1).X
				if firstIP[14] == 0 {
					firstIP[13]++ // 0.(X+1).X.X
					if firstIP[13] == 0 {
						firstIP[12]++ // (X+1).X.X.X
					}
				}
			}
		} else { //IPv6
			var tempIP uint8
			for IPRange.Contains(firstIP) {
				//fmt.Println(firstIP)
				//fmt.Println(firstIP[0], firstIP[1], firstIP[2], firstIP[3], firstIP[4], firstIP[5], firstIP[6], firstIP[7], firstIP[8], firstIP[9], firstIP[10], firstIP[11], firstIP[12], firstIP[13], firstIP[14], firstIP[15])
				firstIP[15] = randipEndWith(255) // 随机 IP 的最后一段
				firstIP[14] = randipEndWith(255) // 随机 IP 的最后一段
				firstIPCopy := make([]byte, len(firstIP))
				copy(firstIPCopy, firstIP)
				firstIPs = append(firstIPs, net.IPAddr{IP: firstIPCopy})
				tempIP = firstIP[13]
				firstIP[13] += randipEndWith(255)
				if firstIP[13] < tempIP {
					tempIP = firstIP[12]
					firstIP[12] += randipEndWith(255)
					if firstIP[12] < tempIP {
						tempIP = firstIP[11]
						firstIP[11] += randipEndWith(255)
						if firstIP[11] < tempIP {
							tempIP = firstIP[10]
							firstIP[10] += randipEndWith(255)
							if firstIP[10] < tempIP {
								tempIP = firstIP[9]
								firstIP[9] += randipEndWith(255)
								if firstIP[9] < tempIP {
									tempIP = firstIP[8]
									firstIP[8] += randipEndWith(255)
									if firstIP[8] < tempIP {
										tempIP = firstIP[7]
										firstIP[7] += randipEndWith(255)
										if firstIP[7] < tempIP {
											tempIP = firstIP[6]
											firstIP[6] += randipEndWith(255)
											if firstIP[6] < tempIP {
												tempIP = firstIP[5]
												firstIP[5] += randipEndWith(255)
												if firstIP[5] < tempIP {
													tempIP = firstIP[4]
													firstIP[4] += randipEndWith(255)
													if firstIP[4] < tempIP {
														tempIP = firstIP[3]
														firstIP[3] += randipEndWith(255)
														if firstIP[3] < tempIP {
															tempIP = firstIP[2]
															firstIP[2] += randipEndWith(255)
															if firstIP[2] < tempIP {
																tempIP = firstIP[1]
																firstIP[1] += randipEndWith(255)
																if firstIP[1] < tempIP {
																	tempIP = firstIP[0]
																	firstIP[0] += randipEndWith(255)
																}
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return firstIPs
}
