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
	OutRegexp         = regexp.MustCompile(`[A-Z]{3}`)
)

// pingReceived pingTotalTime
func (p *Ping) httping(ip *net.IPAddr) (int, time.Duration) {
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
	{
		requ, err := http.NewRequest(http.MethodHead, URL, nil)
		if err != nil {
			return 0, 0
		}
		requ.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.80 Safari/537.36")
		resp, err := hc.Do(requ)
		if err != nil {
			return 0, 0
		}
		defer resp.Body.Close()

		//fmt.Println("IP:", ip, "StatusCode:", resp.StatusCode, resp.Request.URL)
		// 如果未指定的 HTTP 状态码，或指定的状态码不合规，则默认只认为 200、301、302 才算 HTTPing 通过
		if HttpingStatusCode == 0 || HttpingStatusCode < 100 && HttpingStatusCode > 599 {
			if resp.StatusCode != 200 && resp.StatusCode != 301 && resp.StatusCode != 302 {
				return 0, 0
			}
		} else {
			if resp.StatusCode != HttpingStatusCode {
				return 0, 0
			}
		}

		io.Copy(io.Discard, resp.Body)

		// 只有指定了地区才匹配机场三字码
		if HttpingCFColo != "" {
			// 通过头部 Server 值判断是 Cloudflare 还是 AWS CloudFront 并设置 cfRay 为各自的机场三字码完整内容
			cfRay := func() string {
				if resp.Header.Get("Server") == "cloudflare" {
					return resp.Header.Get("CF-RAY") // 示例 cf-ray: 7bd32409eda7b020-SJC
				}
				return resp.Header.Get("x-amz-cf-pop") // 示例 X-Amz-Cf-Pop: SIN52-P1
			}()
			colo := p.getColo(cfRay)
			if colo == "" { // 没有匹配到三字码或不符合指定地区则直接结束该 IP 测试
				return 0, 0
			}
		}

	}

	// 循环测速计算延迟
	success := 0
	var delay time.Duration
	for i := 0; i < PingTimes; i++ {
		requ, err := http.NewRequest(http.MethodHead, URL, nil)
		if err != nil {
			log.Fatal("意外的错误，情报告：", err)
			return 0, 0
		}
		requ.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.80 Safari/537.36")
		if i == PingTimes-1 {
			requ.Header.Set("Connection", "close")
		}
		startTime := time.Now()
		resp, err := hc.Do(requ)
		if err != nil {
			continue
		}
		success++
		io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
		duration := time.Since(startTime)
		delay += duration

	}

	return success, delay

}

func MapColoMap() *sync.Map {
	if HttpingCFColo == "" {
		return nil
	}
	// 将参数指定的地区三字码转为大写并格式化
	colos := strings.Split(strings.ToUpper(HttpingCFColo), ",")
	colomap := &sync.Map{}
	for _, colo := range colos {
		colomap.Store(colo, colo)
	}
	return colomap
}

func (p *Ping) getColo(b string) string {
	if b == "" {
		return ""
	}
	// 正则匹配并返回 机场三字码
	out := OutRegexp.FindString(b)

	if HttpingCFColomap == nil {
		return out
	}
	// 匹配 机场三字码 是否为指定的地区
	_, ok := HttpingCFColomap.Load(out)
	if ok {
		return out
	}

	return ""
}
