package main

import (
	"bufio"
	"log"
	"net"
	"os"
)

func loadFirstIPOfRangeFromFile() []net.IPAddr {
	file, err := os.Open("ip.txt")
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
		firstIP[15]=ipEndWith
		for IPRange.Contains(firstIP) {
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
