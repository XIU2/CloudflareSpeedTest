package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func getCidrHostNum(maskLen int) int {
	cidrIpNum := int(0)
	var i int = int(32 - maskLen - 1)
	for ; i >= 1; i-- {
		cidrIpNum += 1 << i
	}
	return cidrIpNum
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
		firstIP, IPRange, err := net.ParseCIDR(IPString)
		//fmt.Println(firstIP)
		//fmt.Println(IPRange)
		Mask, _ := strconv.Atoi(strings.Split(scanner.Text(), "/")[1])
		MaxIPNum := getCidrHostNum(Mask) - 1
		if MaxIPNum > 253 {
			MaxIPNum = 253
		}
		//fmt.Println(MaxIPNum)
		if err != nil {
			log.Fatal(err)
		}
		if ipv6Mode { // IPv6
			var tempIP uint8
			MaxIPNum = 254
			for IPRange.Contains(firstIP) {
				//fmt.Println(firstIP)
				//fmt.Println(firstIP[0], firstIP[1], firstIP[2], firstIP[3], firstIP[4], firstIP[5], firstIP[6], firstIP[7], firstIP[8], firstIP[9], firstIP[10], firstIP[11], firstIP[12], firstIP[13], firstIP[14], firstIP[15])
				firstIP[15] = randipEndWith(MaxIPNum) // 随机 IP 的最后一段
				firstIP[14] = randipEndWith(MaxIPNum) // 随机 IP 的最后一段
				firstIPCopy := make([]byte, len(firstIP))
				copy(firstIPCopy, firstIP)
				firstIPs = append(firstIPs, net.IPAddr{IP: firstIPCopy})
				tempIP = firstIP[13]
				firstIP[13] += randipEndWith(MaxIPNum)
				if firstIP[13] < tempIP {
					tempIP = firstIP[12]
					firstIP[12] += randipEndWith(MaxIPNum)
					if firstIP[12] < tempIP {
						tempIP = firstIP[11]
						firstIP[11] += randipEndWith(MaxIPNum)
						if firstIP[11] < tempIP {
							tempIP = firstIP[10]
							firstIP[10] += randipEndWith(MaxIPNum)
							if firstIP[10] < tempIP {
								tempIP = firstIP[9]
								firstIP[9] += randipEndWith(MaxIPNum)
								if firstIP[9] < tempIP {
									tempIP = firstIP[8]
									firstIP[8] += randipEndWith(MaxIPNum)
									if firstIP[8] < tempIP {
										tempIP = firstIP[7]
										firstIP[7] += randipEndWith(MaxIPNum)
										if firstIP[7] < tempIP {
											tempIP = firstIP[6]
											firstIP[6] += randipEndWith(MaxIPNum)
											if firstIP[6] < tempIP {
												tempIP = firstIP[5]
												firstIP[5] += randipEndWith(MaxIPNum)
												if firstIP[5] < tempIP {
													tempIP = firstIP[4]
													firstIP[4] += randipEndWith(MaxIPNum)
													if firstIP[4] < tempIP {
														tempIP = firstIP[3]
														firstIP[3] += randipEndWith(MaxIPNum)
														if firstIP[3] < tempIP {
															tempIP = firstIP[2]
															firstIP[2] += randipEndWith(MaxIPNum)
															if firstIP[2] < tempIP {
																tempIP = firstIP[1]
																firstIP[1] += randipEndWith(MaxIPNum)
																if firstIP[1] < tempIP {
																	tempIP = firstIP[0]
																	firstIP[0] += randipEndWith(MaxIPNum)
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
		} else { //IPv4
			for IPRange.Contains(firstIP) {
				//fmt.Println(firstIP)
				//fmt.Println(firstIP[15])
				if allip {
					for i := 1; i < MaxIPNum+2; i++ {
						firstIP[15] = uint8(i) // 随机 IP 的最后一段 0.0.0.X
						//fmt.Println(firstIP)
						firstIPCopy := make([]byte, len(firstIP))
						copy(firstIPCopy, firstIP)
						firstIPs = append(firstIPs, net.IPAddr{IP: firstIPCopy})
					}
				} else {
					if firstIP[15] == 0 {
						firstIP[15] = randipEndWith(MaxIPNum) // 随机 IP 的最后一段 0.0.0.X
					}
					firstIPCopy := make([]byte, len(firstIP))
					copy(firstIPCopy, firstIP)
					firstIPs = append(firstIPs, net.IPAddr{IP: firstIPCopy})
				}
				firstIP[15] = 0
				firstIP[14]++ // 0.0.(X+1).X
				if firstIP[14] == 0 {
					firstIP[13]++ // 0.(X+1).X.X
					if firstIP[13] == 0 {
						firstIP[12]++ // (X+1).X.X.X
					}
				}
			}
		}
	}
	return firstIPs
}
