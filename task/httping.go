package task

import (
	//"crypto/tls"
	//"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	Httping           bool
	HttpingStatusCode int
	HttpingCFColo     string
	HttpingCFColomap  *sync.Map
	ColoRegexp        = regexp.MustCompile(`[A-Z]{3}`)
)

// pingReceived pingTotalTime
func (p *Ping) httping(ip *net.IPAddr) (int, time.Duration, string) {
	hc := http.Client{
		Timeout: time.Second * 2,
		Transport: &http.Transport{
			DialContext: getDialContext(ip),
			//TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 跳过证书验证
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // 阻止重定向
		},
	}

	// 先访问一次获得 HTTP 状态码 及 Cloudflare Colo
	var colo string
	{
		request, err := http.NewRequest(http.MethodHead, URL, nil)
		if err != nil {
			return 0, 0, ""
		}
		request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.80 Safari/537.36")
		response, err := hc.Do(request)
		if err != nil {
			return 0, 0, ""
		}
		defer response.Body.Close()

		//fmt.Println("IP:", ip, "StatusCode:", response.StatusCode, response.Request.URL)
		// 如果未指定的 HTTP 状态码，或指定的状态码不合规，则默认只认为 200、301、302 才算 HTTPing 通过
		if HttpingStatusCode == 0 || HttpingStatusCode < 100 && HttpingStatusCode > 599 {
			if response.StatusCode != 200 && response.StatusCode != 301 && response.StatusCode != 302 {
				return 0, 0, ""
			}
		} else {
			if response.StatusCode != HttpingStatusCode {
				return 0, 0, ""
			}
		}

		io.Copy(io.Discard, response.Body)

		// 通过头部 Server 值判断是 Cloudflare 还是 AWS CloudFront 并设置 cfRay 为各自的机场地区码完整内容
		colo = getHeaderColo(response.Header)

		// 只有指定了地区才匹配机场地区码
		if HttpingCFColo != "" {
			// 判断是否匹配指定的地区码
			colo = p.filterColo(colo)
			if colo == "" { // 没有匹配到地区码或不符合指定地区则直接结束该 IP 测试
				return 0, 0, ""
			}
		}
	}

	// 循环测速计算延迟
	success := 0
	var delay time.Duration
	for i := 0; i < PingTimes; i++ {
		request, err := http.NewRequest(http.MethodHead, URL, nil)
		if err != nil {
			log.Fatal("意外的错误，情报告：", err)
			return 0, 0, ""
		}
		request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.80 Safari/537.36")
		if i == PingTimes-1 {
			request.Header.Set("Connection", "close")
		}
		startTime := time.Now()
		response, err := hc.Do(request)
		if err != nil {
			continue
		}
		success++
		io.Copy(io.Discard, response.Body)
		_ = response.Body.Close()
		duration := time.Since(startTime)
		delay += duration
	}

	return success, delay, colo
}

func MapColoMap() *sync.Map {
	if HttpingCFColo == "" {
		return nil
	}
	// 将 -cfcolo 参数指定的地区地区码转为大写并格式化
	colos := strings.Split(strings.ToUpper(HttpingCFColo), ",")
	colomap := &sync.Map{}
	for _, colo := range colos {
		colomap.Store(colo, colo)
	}
	return colomap
}

// 从响应头中获取 地区码 值
func getHeaderColo(header http.Header) (colo string) {
	// 如果是 Cloudflare 的服务器，则获取 cf-ray 头部
	if header.Get("Server") == "cloudflare" {
		colo = header.Get("cf-ray") // 示例 cf-ray: 7bd32409eda7b020-SJC
	} else { // 反之则默认当成 AWS CloudFront 的服务器获取（如果不是则会获取到空，下面会判断并处理的）
		colo = header.Get("x-amz-cf-pop") // 示例 x-amz-cf-pop: SIN52-P1
	}
	// Fastly CDN 的头部信息，测试地址 https://fastly.jsdelivr.net/gh/XIU2/CloudflareSpeedTest@master/go.mod
	// x-served-by: cache-qpg1275-QPG
	// x-served-by: cache-fra-etou8220141-FRA, cache-hhr-khhr2060043-HHR（最后一个为实际位置）

	// Gcore CDN 的头部信息，测试地址 https://assets.gcore.pro/assets/icons/shield-lock.svg
	// x-id-fe: sg1-hw-edge-gc12
	// x-shard: sg1-shard0-default
	// x-id: sg1-hw-edge-gc2

	// CDN77 的头部信息，测试地址 https://www.cdn77.com
	// Server: CDN77-Turbo
	// X-77-Pop: losangelesUSCA
	// x-77-pop: frankfurtDE
	// x-77-pop: amsterdamNL
	// x-77-pop: singaporeSG

	// 如果没有获取到头部信息，说明不是 Cloudflare 和 AWS CloudFront，则直接返回空字符串
	if colo == "" {
		return ""
	}
	// 正则匹配并返回 机场地区码
	/*matches := ColoRegexp.FindAllString(colo, -1) // 适用于 Fastly 这种有多个地区码的
	if len(matches) == 0 {
		return ""
	}
	return matches[len(matches)-1]*/
	return ColoRegexp.FindString(colo)
}

// 处理地区码
func (p *Ping) filterColo(colo string) string {
	if colo == "" {
		return ""
	}
	// 如果没有指定 -cfcolo 参数，则直接返回
	if HttpingCFColomap == nil {
		return colo
	}
	// 匹配 机场地区码 是否为指定的地区
	_, ok := HttpingCFColomap.Load(colo)
	if ok {
		return colo
	}
	return ""
}
