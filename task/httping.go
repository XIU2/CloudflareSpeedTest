package task

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"CloudflareSpeedTest/utils"
)

var (
	Httping        bool   //是否启用httping
	HttpingTimeout int    //设置超时时间，单位毫秒
	HttpingColo    string //有值代表筛选机场三字码区域
)

var (
	HttpingColomap *sync.Map
	HttpingRequest *http.Request
)

func GetRequestPort(r *http.Request) string {
	port := HttpingRequest.URL.Port()
	if port == "" {
		if HttpingRequest.URL.Scheme == "https" {
			port = "443"
		} else {
			port = "80"
		}
	}
	return port
}

// pingReceived pingTotalTime
func (p *Ping) httping(ip *net.IPAddr) (int, time.Duration) {
	hc := http.Client{
		Timeout: time.Duration(HttpingTimeout) * time.Millisecond,
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				var fullAddress string
				if isIPv4(ip.String()) {
					fullAddress = fmt.Sprintf("%s:%s", ip.String(), GetRequestPort(HttpingRequest))
				} else {
					fullAddress = fmt.Sprintf("[%s]:%s", ip.String(), GetRequestPort(HttpingRequest))
				}
				return (&net.Dialer{}).DialContext(ctx, network, fullAddress)
			},
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	} // #nosec

	traceURL := fmt.Sprintf("%s://www.cloudflare.com:%s/cdn-cgi/trace",
		HttpingRequest.URL.Scheme,
		GetRequestPort(HttpingRequest))

	// for connect and get colo
	{
		requ, err := http.NewRequest(http.MethodGet, traceURL, nil)
		if err != nil {
			return 0, 0
		}
		requ.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.80 Safari/537.36")
		resp, err := hc.Do(requ)
		if err != nil {
			return 0, 0
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return 0, 0
		}

		colo := p.getColo(body)
		if colo == "" {
			return 0, 0
		}

	}

	// for test delay
	success := 0
	var delay time.Duration
	for i := 0; i < PingTimes; i++ {
		requ, err := http.NewRequest(http.MethodGet, traceURL, nil)
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
	if HttpingColo == "" {
		return nil
	}

	colos := strings.Split(HttpingColo, ",")
	colomap := &sync.Map{}
	for _, colo := range colos {
		colomap.Store(colo, colo)
	}
	return colomap
}

func GetRequest() *http.Request {
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Fatal(err)
	}
	return req
}

func (p *Ping) getColo(b []byte) string {
	s := string(b)
	idx := strings.Index(s, "colo=")
	if idx == -1 {
		return ""
	}

	out := s[idx+5 : idx+8]

	utils.ColoAble.Store(out, out)

	if HttpingColomap == nil {
		return out
	}

	_, ok := HttpingColomap.Load(out)
	if ok {
		return out
	}

	return ""
}
