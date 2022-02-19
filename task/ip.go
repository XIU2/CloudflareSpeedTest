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

func InitRandSeed() {
	rand.Seed(time.Now().UnixNano())
}

func randIPEndWith(num byte) byte {
	if num == 0 { // 对于 /32 这种单独的 IP
		return byte(0)
	}
	return byte(rand.Intn(int(num)))
}

type IPRanges struct {
	ips     []*net.IPAddr
	mask    string
	firstIP net.IP
	ipNet   *net.IPNet
}

func newIPRanges() *IPRanges {
	return &IPRanges{
		ips: make([]*net.IPAddr, 0),
	}
}

func (r *IPRanges) fixIP(ip string) string {
	// 如果不含有 '/' 则代表不是 IP 段，而是一个单独的 IP，因此需要加上 /32 /128 子网掩码
	if i := strings.IndexByte(ip, '/'); i < 0 {
		r.mask = "/32"
		if IPv6 {
			r.mask = "/128"
		}
		ip += r.mask
	} else {
		r.mask = ip[i:]
	}
	return ip
}

func (r *IPRanges) parseCIDR(ip string) {
	var err error
	if r.firstIP, r.ipNet, err = net.ParseCIDR(r.fixIP(ip)); err != nil {
		log.Fatalln("ParseCIDR err", err)
	}
}

func (r *IPRanges) appendIPv4(d byte) {
	r.appendIP(net.IPv4(r.firstIP[12], r.firstIP[13], r.firstIP[14], d))
}

func (r *IPRanges) appendIP(ip net.IP) {
	r.ips = append(r.ips, &net.IPAddr{IP: ip})
}

// 返回第四段 ip 的最小值及可用数目
func (r *IPRanges) getIPRange() (minIP, hosts byte) {
	minIP = r.firstIP[15] & r.ipNet.Mask[3] // IP 第四段最小值

	// 根据子网掩码获取主机数量
	m := net.IPv4Mask(255, 255, 255, 255)
	for i, v := range r.ipNet.Mask {
		m[i] ^= v
	}
	total, _ := strconv.ParseInt(m.String(), 16, 32) // 总可用 IP 数
	if total > 255 {                                 // 矫正 第四段 可用 IP 数
		hosts = 255
		return
	}
	hosts = byte(total)
	return
}

func (r *IPRanges) chooseIPv4() {
	minIP, hosts := r.getIPRange()
	for r.ipNet.Contains(r.firstIP) {
		if TestAll { // 如果是测速全部 IP
			for i := 0; i <= int(hosts); i++ { // 遍历 IP 最后一段最小值到最大值
				r.appendIPv4(byte(i) + minIP)
			}
		} else { // 随机 IP 的最后一段 0.0.0.X
			r.appendIPv4(minIP + randIPEndWith(hosts))
		}
		r.firstIP[14]++ // 0.0.(X+1).X
		if r.firstIP[14] == 0 {
			r.firstIP[13]++ // 0.(X+1).X.X
			if r.firstIP[13] == 0 {
				r.firstIP[12]++ // (X+1).X.X.X
			}
		}
	}
}

func (r *IPRanges) chooseIPv6() {
	var tempIP uint8
	for r.ipNet.Contains(r.firstIP) {
		//fmt.Println(firstIP)
		//fmt.Println(firstIP[0], firstIP[1], firstIP[2], firstIP[3], firstIP[4], firstIP[5], firstIP[6], firstIP[7], firstIP[8], firstIP[9], firstIP[10], firstIP[11], firstIP[12], firstIP[13], firstIP[14], firstIP[15])
		if r.mask != "/128" {
			r.firstIP[15] = randIPEndWith(255) // 随机 IP 的最后一段
			r.firstIP[14] = randIPEndWith(255) // 随机 IP 的最后一段
		}
		targetIP := make([]byte, len(r.firstIP))
		copy(targetIP, r.firstIP)
		r.appendIP(targetIP)
		for i := 13; i >= 0; i-- {
			tempIP = r.firstIP[i]
			r.firstIP[i] += randIPEndWith(255)
			if r.firstIP[i] >= tempIP {
				break
			}
		}
	}
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
	ranges := newIPRanges()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ranges.parseCIDR(scanner.Text())
		if IPv6 {
			ranges.chooseIPv6()
			continue
		}
		ranges.chooseIPv4()
	}
	return ranges.ips
}
