package main

import (
	"bufio"
	"log"
	"net"
	"os"
)

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
		if err != nil {
			log.Fatal(err)
		}
		if ipv6Mode { // IPv6
			var tempIP uint8
			for IPRange.Contains(firstIP) {
				//fmt.Println(firstIP)
				//fmt.Println(firstIP[0], firstIP[1], firstIP[2], firstIP[3], firstIP[4], firstIP[5], firstIP[6], firstIP[7], firstIP[8], firstIP[9], firstIP[10], firstIP[11], firstIP[12], firstIP[13], firstIP[14], firstIP[15])
				firstIP[15] = randipEndWith() // 随机 IP 的最后一段
				firstIP[14] = randipEndWith() // 随机 IP 的最后一段
				firstIPCopy := make([]byte, len(firstIP))
				copy(firstIPCopy, firstIP)
				firstIPs = append(firstIPs, net.IPAddr{IP: firstIPCopy})
				tempIP = firstIP[13]
				firstIP[13] += randipEndWith()
				if firstIP[13] < tempIP {
					tempIP = firstIP[12]
					firstIP[12] += randipEndWith()
					if firstIP[12] < tempIP {
						tempIP = firstIP[11]
						firstIP[11] += randipEndWith()
						if firstIP[11] < tempIP {
							tempIP = firstIP[10]
							firstIP[10] += randipEndWith()
							if firstIP[10] < tempIP {
								tempIP = firstIP[9]
								firstIP[9] += randipEndWith()
								if firstIP[9] < tempIP {
									tempIP = firstIP[8]
									firstIP[8] += randipEndWith()
									if firstIP[8] < tempIP {
										tempIP = firstIP[7]
										firstIP[7] += randipEndWith()
										if firstIP[7] < tempIP {
											tempIP = firstIP[6]
											firstIP[6] += randipEndWith()
											if firstIP[6] < tempIP {
												tempIP = firstIP[5]
												firstIP[5] += randipEndWith()
												if firstIP[5] < tempIP {
													tempIP = firstIP[4]
													firstIP[4] += randipEndWith()
													if firstIP[4] < tempIP {
														tempIP = firstIP[3]
														firstIP[3] += randipEndWith()
														if firstIP[3] < tempIP {
															tempIP = firstIP[2]
															firstIP[2] += randipEndWith()
															if firstIP[2] < tempIP {
																tempIP = firstIP[1]
																firstIP[1] += randipEndWith()
																if firstIP[1] < tempIP {
																	tempIP = firstIP[0]
																	firstIP[0] += randipEndWith()
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
				firstIP[15] = randipEndWith() // 随机 IP 的最后一段 0.0.0.X
				firstIPCopy := make([]byte, len(firstIP))
				copy(firstIPCopy, firstIP)
				firstIPs = append(firstIPs, net.IPAddr{IP: firstIPCopy})
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
