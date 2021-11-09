package tcp

import (
	"fmt"
	"net"
	"sync"
	"time"

	"CloudflareSpeedTest/utils"
)

const tcpConnectTimeout = time.Second * 1

type Ping struct {
	wg              *sync.WaitGroup
	m               *sync.Mutex
	ips             []net.IPAddr
	isIPv6          bool
	tcpPort         int
	pingCount       int
	csv             []utils.CloudflareIPData
	control         chan bool
	progressHandler func(e utils.ProgressEvent)
}

func NewPing(ips []net.IPAddr, port, pingTime int, ipv6 bool) *Ping {
	return &Ping{
		wg:        &sync.WaitGroup{},
		m:         &sync.Mutex{},
		ips:       ips,
		isIPv6:    ipv6,
		tcpPort:   port,
		pingCount: pingTime,
		csv:       make([]utils.CloudflareIPData, 0),
		control:   make(chan bool),
	}
}

func (p *Ping) Run() {
	for _, ip := range p.ips {
		p.wg.Add(1)
		p.control <- false
		go p.start(ip)
	}
}

func (p *Ping) start(ip net.IPAddr) {
	defer p.wg.Done()
	if ok, data := p.tcpingHandler(ip, nil); ok {
		p.appendIPData(data)
	}
	<-p.control
}

func (p *Ping) appendIPData(data *utils.PingData) {
	p.m.Lock()
	defer p.m.Unlock()
	p.csv = append(p.csv, utils.CloudflareIPData{
		PingData: data,
	})
}

//bool connectionSucceed float32 time
func (p *Ping) tcping(ip net.IPAddr) (bool, time.Duration) {
	startTime := time.Now()
	fullAddress := fmt.Sprintf("%s:%d", ip.String(), p.tcpPort)
	//fmt.Println(ip.String())
	if p.isIPv6 { // IPv6 需要加上 []
		fullAddress = fmt.Sprintf("[%s]:%d", ip.String(), p.tcpPort)
	}
	conn, err := net.DialTimeout("tcp", fullAddress, tcpConnectTimeout)
	if err != nil {
		return false, 0
	}
	defer conn.Close()
	duration := time.Since(startTime)
	return true, duration
}

//pingReceived pingTotalTime
func (p *Ping) checkConnection(ip net.IPAddr) (pingRecv int, pingTime time.Duration) {
	for i := 0; i < p.pingCount; i++ {
		if pingSucceed, pingTimeCurrent := p.tcping(ip); pingSucceed {
			pingRecv++
			pingTime += pingTimeCurrent
		}
	}
	return
}

//return Success packetRecv averagePingTime specificIPAddr
func (p *Ping) tcpingHandler(ip net.IPAddr, progressHandler func(e utils.ProgressEvent)) (bool, *utils.PingData) {
	ipCanConnect := false
	pingRecv := 0
	var pingTime time.Duration
	for !ipCanConnect {
		pingRecvCurrent, pingTimeCurrent := p.checkConnection(ip)
		if pingRecvCurrent != 0 {
			ipCanConnect = true
			pingRecv = pingRecvCurrent
			pingTime = pingTimeCurrent
		} else {
			ip.IP[15]++
			if ip.IP[15] == 0 {
				break
			}
			break
		}
	}
	if !ipCanConnect {
		progressHandler(utils.NoAvailableIPFound)
		return false, nil
	}
	progressHandler(utils.AvailableIPFound)
	for i := 0; i < p.pingCount; i++ {
		pingSuccess, pingTimeCurrent := p.tcping(ip)
		progressHandler(utils.NormalPing)
		if pingSuccess {
			pingRecv++
			pingTime += pingTimeCurrent
		}
	}
	return true, &utils.PingData{
		IP:       ip,
		Count:    p.pingCount,
		Received: pingRecv,
		Delay:    pingTime / time.Duration(pingRecv),
	}
}
