package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/VividCortex/ewma"
)

//bool connectionSucceed float32 time
func tcping(ip net.IPAddr, tcpPort int) (bool, time.Duration) {
	startTime := time.Now()
	fullAddress := fmt.Sprintf("%s:%d", ip.String(), tcpPort)
	//fmt.Println(ip.String())
	if ipv6Mode { // IPv6 需要加上 []
		fullAddress = fmt.Sprintf("[%s]:%d", ip.String(), tcpPort)
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
func checkConnection(ip net.IPAddr, tcpPort int) (pingRecv int, pingTime time.Duration) {
	for i := 1; i <= failTime; i++ {
		pingSucceed, pingTimeCurrent := tcping(ip, tcpPort)
		if pingSucceed {
			pingRecv++
			pingTime += pingTimeCurrent
		}
	}
	return
}

//return Success packetRecv averagePingTime specificIPAddr
func tcpingHandler(ip net.IPAddr, tcpPort, pingCount int, progressHandler func(e progressEvent)) (bool, int, time.Duration, net.IPAddr) {
	ipCanConnect := false
	pingRecv := 0
	var pingTime time.Duration
	for !ipCanConnect {
		pingRecvCurrent, pingTimeCurrent := checkConnection(ip, tcpPort)
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
		progressHandler(NoAvailableIPFound)
		return false, 0, 0, net.IPAddr{}
	}
	progressHandler(AvailableIPFound)
	for i := failTime; i < pingCount; i++ {
		fmt.Println("failTime", failTime)
		pingSuccess, pingTimeCurrent := tcping(ip, tcpPort)
		progressHandler(NormalPing)
		if pingSuccess {
			pingRecv++
			pingTime += pingTimeCurrent
		}
	}
	return true, pingRecv, pingTime / time.Duration(pingRecv), ip
}

func tcpingGoroutine(wg *sync.WaitGroup, mutex *sync.Mutex, ip net.IPAddr, tcpPort int, pingCount int, csv *[]CloudflareIPData, control chan bool, progressHandler func(e progressEvent)) {
	defer wg.Done()
	// fmt.Println(ip.String())
	success, pingRecv, pingTimeAvg, currentIP := tcpingHandler(ip, tcpPort, pingCount, progressHandler)
	if success {
		mutex.Lock()
		var cfdata CloudflareIPData
		cfdata.ip = currentIP
		cfdata.pingReceived = pingRecv
		cfdata.pingTime = pingTimeAvg
		cfdata.pingCount = pingCount
		*csv = append(*csv, cfdata)
		mutex.Unlock()
	}
	// <-control
}

func GetDialContextByAddr(fakeSourceAddr string) func(ctx context.Context, network, address string) (net.Conn, error) {
	return func(ctx context.Context, network, address string) (net.Conn, error) {
		c, e := (&net.Dialer{}).DialContext(ctx, network, fakeSourceAddr)
		return c, e
	}
}

//bool : can download,float32 downloadSpeed
func DownloadSpeedHandler(ip net.IPAddr) (bool, float32) {
	var client = http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       downloadTestTime,
	}
	var fullAddress string
	if ipv6Mode { // IPv6 需要加上 []
		fullAddress = "[" + ip.String() + "]:443"
	} else {
		fullAddress = ip.String() + ":443"
	}
	client.Transport = &http.Transport{
		DialContext: GetDialContextByAddr(fullAddress),
	}
	response, err := client.Get(url)
	if err != nil {
		return false, 0
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return false, 0
	}
	timeStart := time.Now()
	timeEnd := timeStart.Add(downloadTestTime)

	contentLength := response.ContentLength
	buffer := make([]byte, downloadBufferSize)

	var (
		contentRead     int64 = 0
		timeSlice             = downloadTestTime / 100
		timeCounter           = 1
		lastContentRead int64 = 0
	)

	var nextTime = timeStart.Add(timeSlice * time.Duration(timeCounter))
	e := ewma.NewMovingAverage()

	for contentLength != contentRead {
		var currentTime = time.Now()
		if currentTime.After(nextTime) {
			timeCounter++
			nextTime = timeStart.Add(timeSlice * time.Duration(timeCounter))
			e.Add(float64(contentRead - lastContentRead))
			lastContentRead = contentRead
		}
		if currentTime.After(timeEnd) {
			break
		}
		bufferRead, err := response.Body.Read(buffer)
		contentRead += int64(bufferRead)
		if err != nil {
			if err != io.EOF {
				break
			}
			e.Add(float64(contentRead-lastContentRead) / (float64(nextTime.Sub(currentTime)) / float64(timeSlice)))
		}
	}
	return true, float32(e.Value()) / (float32(downloadTestTime.Seconds()) / 120)

}
