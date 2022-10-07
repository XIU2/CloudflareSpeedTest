package task

import (
	"bufio"
	"crypto/rand"
	"log"
	"math/big"
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

		var addrs []netip.Addr

		switch {
		case prefix.IsSingleIP():
			// if /32 or /128, append to the array directly
			addrs = append(addrs, prefix.Addr())

		case IPv6 || !TestAll:
			// here, choose an random address

			// the network address
			netAddr := prefix.Addr()
			// how many bits can hosts take
			hostBits := netAddr.BitLen() - prefix.Bits()
			// the last host address in this subnet + 1
			bigHostAddrMax := new(big.Int).Lsh(big.NewInt(1), uint(hostBits))
			// take one address
			bigHostAddr, err := rand.Int(rand.Reader, bigHostAddrMax)
			if err != nil {
				panic("Failed to choose random address!")
			}

			// convert netip.Addr to big.Int
			bigNetAddr := new(big.Int).SetBytes(netAddr.AsSlice())
			// network address + host part
			bigAddr := new(big.Int).Or(bigNetAddr, bigHostAddr)
			// convert big.Int to netip.Addr
			addr, _ := netip.AddrFromSlice(bigAddr.Bytes())

			addrs = append(addrs, addr)

		default:
			// we need all the possible addresses
			for addr := prefix.Addr(); prefix.Contains(addr); addr = addr.Next() {
				addrs = append(addrs, addr)
			}
			addrs = addrs[1 : len(addrs)-1]
		}

		// add them to the original "ips" array
		for _, v := range addrs {
			ips = append(ips, &net.IPAddr{IP: net.IP(v.AsSlice())})
		}
	}

	return ips
}
