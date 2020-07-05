package main

import (
	"net"
	"strconv"
	"sync"
	"time"
)

const defaultTcpPort = 443
const tcpConnectTimeout = time.Second * 1
const failTime = 4

//bool connectionSucceed float64 time
func tcping(ip net.IPAddr) (bool, float64) {
	startTime := time.Now()
	conn, err := net.DialTimeout("tcp", ip.String()+":"+strconv.Itoa(defaultTcpPort), tcpConnectTimeout)
	if err != nil {
		return false, 0
	} else {
		var endTime = time.Since(startTime)
		var duration = float64(endTime.Microseconds()) / 1000.0
		_ = conn.Close()
		return true, duration
	}
}

//pingReceived pingTotalTime
func checkConnection(ip net.IPAddr) (int, float64) {
	pingRecv := 0
	pingTime := 0.0
	for i := 1; i <= failTime; i++ {
		pingSucceed, pingTimeCurrent := tcping(ip)
		if pingSucceed {
			pingRecv++
			pingTime += pingTimeCurrent
		}
	}
	return pingRecv, pingTime
}

//return Success packetRecv averagePingTime specificIPAddr
func tcpingHandler(ip net.IPAddr, pingCount int, progressHandler func(e progressEvent)) (bool, int, float64, net.IPAddr) {
	ipCanConnect := false
	pingRecv := 0
	pingTime := 0.0
	for !ipCanConnect {
		pingRecvCurrent, pingTimeCurrent := checkConnection(ip)
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
	if ipCanConnect {
		progressHandler(AvailableIPFound)
		for i := failTime; i < pingCount; i++ {
			pingSuccess, pingTimeCurrent := tcping(ip)
			progressHandler(NormalPing)
			if pingSuccess {
				pingRecv++
				pingTime += pingTimeCurrent
			}
		}
		return true, pingRecv, pingTime / float64(pingRecv), ip
	} else {
		progressHandler(NoAvailableIPFound)
		return false, 0, 0, net.IPAddr{}
	}
}

func tcpingGoroutine(wg *sync.WaitGroup, mutex *sync.Mutex, ip net.IPAddr, pingCount int, csv *[][]string, control chan bool, progressHandler func(e progressEvent)) {
	defer wg.Done()
	success, pingRecv, pingTimeAvg, currentIP := tcpingHandler(ip, pingCount, progressHandler)
	if success {
		mutex.Lock()
		*csv = append(*csv, []string{currentIP.String(), strconv.Itoa(pingRecv), strconv.FormatFloat(pingTimeAvg, 'f', 4, 64)})
		mutex.Unlock()
	}
	<-control
}
