package task

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"CloudflareSpeedTest/utils"
)

const (
	tcpConnectTimeout = time.Second * 1
	maxRoutine        = 1000
	defaultRoutines   = 200
	defaultPort       = 443
	defaultPingTimes  = 4
)

var (
	Routines      = defaultRoutines
	TCPPort   int = defaultPort
	PingTimes int = defaultPingTimes
	Colo      string
)

type Ping struct {
	wg      *sync.WaitGroup
	m       *sync.Mutex
	ips     []*net.IPAddr
	csv     utils.PingDelaySet
	control chan bool
	bar     *utils.Bar
	colomap *sync.Map
	request *http.Request
}

func checkPingDefault() {
	if Routines <= 0 {
		Routines = defaultRoutines
	}
	if TCPPort <= 0 || TCPPort >= 65535 {
		TCPPort = defaultPort
	}
	if PingTimes <= 0 {
		PingTimes = defaultPingTimes
	}
}

func NewPing() *Ping {
	checkPingDefault()
	ips := loadIPRanges()
	return &Ping{
		wg:      &sync.WaitGroup{},
		m:       &sync.Mutex{},
		ips:     ips,
		csv:     make(utils.PingDelaySet, 0),
		control: make(chan bool, Routines),
		bar:     utils.NewBar(len(ips)),
		colomap: mapColoMap(),
		request: getRequest(),
	}
}

func (p *Ping) Run() utils.PingDelaySet {
	if len(p.ips) == 0 {
		return p.csv
	}
	fmt.Printf("开始延迟测速（模式：TCP，端口：%d，平均延迟上限：%v ms，平均延迟下限：%v ms)\n", TCPPort, utils.InputMaxDelay.Milliseconds(), utils.InputMinDelay.Milliseconds())
	for _, ip := range p.ips {
		p.wg.Add(1)
		p.control <- false
		go p.start(ip)
	}
	p.wg.Wait()
	p.bar.Done()
	sort.Sort(p.csv)
	return p.csv
}

func (p *Ping) start(ip *net.IPAddr) {
	defer p.wg.Done()
	p.tcpingHandler(ip)
	<-p.control
}

// bool connectionSucceed float32 time
func (p *Ping) tcping(ip *net.IPAddr) (bool, time.Duration) {
	startTime := time.Now()
	var fullAddress string
	if isIPv4(ip.String()) {
		fullAddress = fmt.Sprintf("%s:%d", ip.String(), TCPPort)
	} else {
		fullAddress = fmt.Sprintf("[%s]:%d", ip.String(), TCPPort)
	}
	conn, err := net.DialTimeout("tcp", fullAddress, tcpConnectTimeout)
	if err != nil {
		return false, 0
	}
	defer conn.Close()
	duration := time.Since(startTime)
	return true, duration
}

// pingReceived pingTotalTime
func (p *Ping) checkConnection(ip *net.IPAddr) (recv int, totalDelay time.Duration) {
	if p.colomap != nil {
		recv, totalDelay = p.httpping(ip)
		return
	}
	for i := 0; i < PingTimes; i++ {
		if ok, delay := p.tcping(ip); ok {
			recv++
			totalDelay += delay
		}
	}
	return
}

func (p *Ping) appendIPData(data *utils.PingData) {
	p.m.Lock()
	defer p.m.Unlock()
	p.csv = append(p.csv, utils.CloudflareIPData{
		PingData: data,
	})
}

// handle tcping
func (p *Ping) tcpingHandler(ip *net.IPAddr) {
	recv, totalDlay := p.checkConnection(ip)
	p.bar.Grow(1)
	if recv == 0 {
		return
	}
	data := &utils.PingData{
		IP:       ip,
		Sended:   PingTimes,
		Received: recv,
		Delay:    totalDlay / time.Duration(recv),
	}
	p.appendIPData(data)
}

// pingReceived pingTotalTime
func (p *Ping) httpping(ip *net.IPAddr) (int, time.Duration) {

	hc := http.Client{
		Timeout: Timeout,
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				var fullAddress string
				if isIPv4(ip.String()) {
					fullAddress = fmt.Sprintf("%s:%d", ip.String(), TCPPort)
				} else {
					fullAddress = fmt.Sprintf("[%s]:%d", ip.String(), TCPPort)
				}
				return (&net.Dialer{}).DialContext(ctx, network, fullAddress)
			},
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	} // #nosec

	traceURL := fmt.Sprintf("%s://www.cloudflare.com:%s/cdn-cgi/trace",
		p.request.URL.Scheme,
		p.request.URL.Port())

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

	startTime := time.Now()
	// for test delay
	for i := 0; i < PingTimes; i++ {
		requ, err := http.NewRequest(http.MethodGet, traceURL, nil)
		if err != nil {
			return 0, 0
		}
		requ.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.80 Safari/537.36")
		if i == PingTimes-1 {
			requ.Header.Set("Connection", "close")
		}
		resp, err := hc.Do(requ)
		if err != nil {
			return 0, 0
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
	duration := time.Since(startTime)
	return PingTimes, duration

}

func mapColoMap() *sync.Map {
	if Colo == "" {
		return nil
	}

	colos := strings.Split(Colo, ",")
	colomap := &sync.Map{}
	for _, colo := range colos {
		colomap.Store(colo, colo)
	}
	return colomap
}

func getRequest() *http.Request {
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

	_, ok := p.colomap.Load(out)
	if ok {
		return out
	}

	return ""
}
