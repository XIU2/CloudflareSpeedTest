package main

import (
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"sort"
	"sync"
	"time"
)

func handleUserInput() {
	fmt.Println("请输入扫描协程数(数字越大越快,默认400):")
	fmt.Scanln(&pingRoutine)
	if pingRoutine <= 0 {
		pingRoutine = 400
	}
	fmt.Println("请输入tcping次数(默认10):")
	fmt.Scanln(&pingTime)
	if pingTime <= 0 {
		pingTime = 10
	}
	fmt.Println("请输入要测试的下载节点个数(默认10):")
	fmt.Scanln(&downloadTestCount)
	if downloadTestCount <= 0 {
		downloadTestCount = 10
	}
	fmt.Println("请输入下载测试时间(默认10,单位为秒):")
	var downloadSecond int64
	fmt.Scanln(&downloadSecond)
	if downloadSecond <= 0 {
		downloadSecond = 10
	}
	downloadTestTime = time.Duration(downloadSecond) * time.Second
}

func main() {
	initipEndWith()
	handleUserInput()
	ips := loadFirstIPOfRangeFromFile()
	pingCount := len(ips) * pingTime
	bar := pb.StartNew(pingCount)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var data = make([]CloudflareIPData, 0)

	fmt.Println("开始tcping")

	control := make(chan bool, pingRoutine)
	for _, ip := range ips {
		wg.Add(1)
		control <- false
		handleProgress := handleProgressGenerator(bar)
		go tcpingGoroutine(&wg, &mu, ip, pingTime, &data, control, handleProgress)
	}
	wg.Wait()
	bar.Finish()
	bar = pb.StartNew(downloadTestCount)
	sort.Sort(CloudflareIPDataSet(data))
	fmt.Println("开始下载测速")
	for i := 0; i < downloadTestCount; i++ {
		_, speed := DownloadSpeedHandler(data[i].ip)
		data[i].downloadSpeed = speed
		bar.Add(1)
	}
	bar.Finish()
	ExportCsv("./result.csv", data)
}
