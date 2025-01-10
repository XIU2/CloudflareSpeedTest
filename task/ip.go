package task

import (
	"io"
	"log"
	random "math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const defaultInputFile = "ip.txt"
const defaultRemoteURL = "https://www.cloudflare.com/ips-v4" // ipv6: https://www.cloudflare.com/ips-v6

var (
	// TestAll test all ip
	TestAll = false
	// UseRemoteURL use remote url rather than local file
	UseRemoteURL = false

	// IPFile is the filename of IP Ranges
	IPFile = defaultInputFile
	// IPRemoteURL is the remote url of IP Ranges
	IPRemoteURL = defaultRemoteURL
	IPText      string

	rand *random.Rand
)

func InitRandSeed() {
	rand = random.New(random.NewSource(time.Now().UnixNano()))
}

func isIPv4(ip string) bool {
	return strings.Contains(ip, ".")
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

// 如果是单独 IP 则加上子网掩码，反之则获取子网掩码(r.mask)
func (r *IPRanges) fixIP(ip string) string {
	// 如果不含有 '/' 则代表不是 IP 段，而是一个单独的 IP，因此需要加上 /32 /128 子网掩码
	if i := strings.IndexByte(ip, '/'); i < 0 {
		if isIPv4(ip) {
			r.mask = "/32"
		} else {
			r.mask = "/128"
		}
		ip += r.mask
	} else {
		r.mask = ip[i:]
	}
	return ip
}

// 解析 IP 段，获得 IP、IP 范围、子网掩码
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
	if r.mask == "/32" { // 单个 IP 则无需随机，直接加入自身即可
		r.appendIP(r.firstIP)
	} else {
		minIP, hosts := r.getIPRange()    // 返回第四段 IP 的最小值及可用数目
		for r.ipNet.Contains(r.firstIP) { // 只要该 IP 没有超出 IP 网段范围，就继续循环随机
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
}

func (r *IPRanges) chooseIPv6() {
	if r.mask == "/128" { // 单个 IP 则无需随机，直接加入自身即可
		r.appendIP(r.firstIP)
	} else {
		var tempIP uint8                  // 临时变量，用于记录前一位的值
		for r.ipNet.Contains(r.firstIP) { // 只要该 IP 没有超出 IP 网段范围，就继续循环随机
			r.firstIP[15] = randIPEndWith(255) // 随机 IP 的最后一段
			r.firstIP[14] = randIPEndWith(255) // 随机 IP 的最后一段

			targetIP := make([]byte, len(r.firstIP))
			copy(targetIP, r.firstIP)
			r.appendIP(targetIP) // 加入 IP 地址池

			for i := 13; i >= 0; i-- { // 从倒数第三位开始往前随机
				tempIP = r.firstIP[i]              // 保存前一位的值
				r.firstIP[i] += randIPEndWith(255) // 随机 0~255，加到当前位上
				if r.firstIP[i] >= tempIP {        // 如果当前位的值大于等于前一位的值，说明随机成功了，可以退出该循环
					break
				}
			}
		}
	}
}

func loadFromIPSegment(ranges *IPRanges, ipSegment string) {
	ipSegment = strings.TrimSpace(ipSegment) // 去除首尾的空白字符（空格、制表符、换行符等）
	if ipSegment == "" {                     // 跳过空行
		return
	}
	ranges.parseCIDR(ipSegment) // 解析 IP 段，获得 IP、IP 范围、子网掩码
	if isIPv4(ipSegment) {      // 生成要测速的所有 IPv4 / IPv6 地址（单个/随机/全部）
		ranges.chooseIPv4()
	} else {
		ranges.chooseIPv6()
	}
}

func ipSegmentsFromFile(file string) []string {
	bs, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	return strings.Split(string(bs), "\n")
}

func ipSegmentsFromURL(url string) []string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return strings.Split(string(bs), "\n")
}

func loadIPRanges() []*net.IPAddr {
	ranges := newIPRanges()
	ipSegments := []string{}

	if IPText != "" { // 从参数中获取 IP 段数据
		ipSegments = strings.Split(IPText, ",") // 以逗号分隔为数组并循环遍历
	} else if UseRemoteURL { // 从远程 URL 获取 IP 段数据
		if IPRemoteURL == "" {
			IPRemoteURL = defaultRemoteURL
		}
		ipSegments = ipSegmentsFromURL(IPRemoteURL)
	} else { // 从文件中获取 IP 段数据
		if IPFile == "" {
			IPFile = defaultInputFile
		}
		ipSegments = ipSegmentsFromFile(IPFile)
	}

	for _, segment := range ipSegments {
		loadFromIPSegment(ranges, segment)
	}
	return ranges.ips
}
