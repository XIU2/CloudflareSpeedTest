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
		for IPRange.Contains(firstIP) {
			randipEndWith() // 随机 IP 的最后一段
			firstIP[15] = ipEndWith
			firstIPCopy := make([]byte, len(firstIP))
			copy(firstIPCopy, firstIP)
			firstIPs = append(firstIPs, net.IPAddr{IP: firstIPCopy})
			firstIP[14]++
			if firstIP[14] == 0 {
				firstIP[13]++
				if firstIP[13] == 0 {
					firstIP[12]++
				}
			}
		}
	}
	return firstIPs
}
