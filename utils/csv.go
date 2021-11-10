package utils

import (
	"encoding/csv"
	"log"
	"net"
	"os"
	"sort"
	"strconv"
	"time"
)

var (
	MaxDelay = 9999 * time.Millisecond
	MinDelay = time.Duration(0)

	InputMaxDelay = MaxDelay
	InputMinDelay = MinDelay
)

type PingData struct {
	IP       net.IPAddr
	Sended   int
	Received int
	Delay    time.Duration
}

type CloudflareIPData struct {
	*PingData
	recvRate      float32
	downloadSpeed float32
}

func (cf *CloudflareIPData) getRecvRate() float32 {
	if cf.recvRate == 0 {
		pingLost := cf.Sended - cf.Received
		cf.recvRate = float32(pingLost) / float32(cf.Sended)
	}
	return cf.recvRate
}

func (cf *CloudflareIPData) toString() []string {
	result := make([]string, 6)
	result[0] = cf.IP.String()
	result[1] = strconv.Itoa(cf.Sended)
	result[2] = strconv.Itoa(cf.Received)
	result[3] = strconv.FormatFloat(float64(cf.getRecvRate()), 'f', 2, 32)
	result[4] = cf.Delay.String()
	result[5] = strconv.FormatFloat(float64(cf.downloadSpeed)/1024/1024, 'f', 2, 32)
	return result
}

func ExportCsv(filePath string, data []CloudflareIPData) {
	fp, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("创建文件[%s]失败：%v", filePath, err)
		return
	}
	defer fp.Close()
	w := csv.NewWriter(fp) //创建一个新的写入文件流
	w.Write([]string{"IP 地址", "已发送", "已接收", "丢包率", "平均延迟", "下载速度 (MB/s)"})
	w.WriteAll(convertToString(data))
	w.Flush()
}

func convertToString(data []CloudflareIPData) [][]string {
	result := make([][]string, 0)
	for _, v := range data {
		result = append(result, v.toString())
	}
	return result
}

type PingDelaySet []CloudflareIPData

func (s PingDelaySet) FilterDelay() (data PingDelaySet) {
	sort.Sort(s)
	if InputMaxDelay >= MaxDelay || InputMinDelay <= MinDelay {
		return s
	}
	for _, v := range s {
		if v.Delay > MaxDelay { // 平均延迟上限
			break
		}
		if v.Delay <= MinDelay { // 平均延迟下限
			continue
		}
		data = append(data, v) // 延迟满足条件时，添加到新数组中
	}
	return
}

func (s PingDelaySet) Len() int {
	return len(s)
}

func (s PingDelaySet) Less(i, j int) bool {
	iRate, jRate := s[i].getRecvRate(), s[j].getRecvRate()
	if iRate != jRate {
		return iRate < jRate
	}
	return s[i].Delay < s[j].Delay
}

func (s PingDelaySet) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// 下载速度排序
type DownloadSpeedSet []CloudflareIPData

func (s DownloadSpeedSet) Len() int {
	return len(s)
}

func (s DownloadSpeedSet) Less(i, j int) bool {
	return s[i].downloadSpeed > s[j].downloadSpeed
}

func (s DownloadSpeedSet) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
