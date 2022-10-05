package task

import (
	"bufio"
	"log"
	"math/rand"
	"net"
	"net/netip"
	"os"
	"regexp"
)

var (
	// Use ipv6 in the test
	IPv6 = false
	// If true, test all the addresses
	TestAll = false
	// The filename of ip range data
	IPFile string
)

func rangeNormalize(s string) string {
	noMask := !regexp.MustCompile(`.*\/[0-9]+`).MatchString(s)
	if noMask {
		addr := netip.MustParseAddr(s)
		switch {
		case addr.Is4():
			return s + "/32"
		case addr.Is6():
			return s + "/128"
		}
	}
	return s
}

func loadIPRanges() []*net.IPAddr {

	var ips []*net.IPAddr

	// try to open data file
	file, err := os.Open(IPFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// load in file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// normalize the expression
		line := rangeNormalize(scanner.Text())

		// parse the network part
		prefix := netip.MustParsePrefix(line)

		// ensure we are using the standard form
		prefix = prefix.Masked()

		// get all possible address
		var addrs []netip.Addr
		if prefix.IsSingleIP() {
			for addr := prefix.Addr(); prefix.Contains(addr); addr = addr.Next() {
				addrs = append(addrs, addr)
			}
			addrs = addrs[1 : len(addrs)-1]
		} else {
			addrs = append(addrs, prefix.Addr())
		}

		// choose one addr if we don't need all
		if !TestAll {
			addrs = []netip.Addr{addrs[rand.Intn(len(addrs))]}
		}

		// add them to the original "ips" array
		for _, v := range addrs {
			ips = append(ips, &net.IPAddr{IP: net.IP(v.AsSlice())})
		}
	}

	return ips
}
