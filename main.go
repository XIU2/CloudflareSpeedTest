package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/cheggaaa/pb/v3"
)

var version string

func init() {
	var downloadSecond int64
	var printVersion bool
	const help = `CloudflareSpeedTest
测试 Cloudflare CDN 所有 IP 的延迟和速度，获取最佳 IP！
https://github.com/XIU2/CloudflareSpeedTest

参数：
    -n 500
        测速线程数量；请勿超过1000 (默认 500)
    -t 4
        延迟测速次数；单个 IP (默认 4)
    -dn 20
        下载测速数量；延迟测速后，从最低延迟起测试下载速度的数量，请勿太多 (默认 20)
    -dt 10
        下载测试时间；单个 IP 测速最长时间，单位：秒 (默认 10)
    -v
        打印程序版本
    -h
        打印帮助说明

示例：
    Windows：CloudflareST.exe -n 800 -t 4 -dn 20 -dt 10
    Linux：CloudflareST -n 800 -t 4 -dn 20 -dt 10
`

	pingRoutine = *flag.Int("n", 500, "测速线程数量；请勿超过1000")
	pingTime = *flag.Int("t", 4, "延迟测速次数；单个 IP")
	downloadTestCount = *flag.Int("dn", 20, "下载测速数量；延迟测速后，从最低延迟起测试下载速度的数量，请勿太多")
	flag.Int64Var(&downloadSecond, "dt", 10, "下载测速时间；单个 IP 测速最长时间，单位：秒")
	flag.BoolVar(&printVersion, "v", false, "打印程序版本")

	downloadTestTime = time.Duration(downloadSecond) * time.Second

	flag.Usage = func() { fmt.Print(help) }
	flag.Parse()
	if printVersion {
		println(version)
		os.Exit(0)
	}
}

func main() {
	initipEndWith()
	ips := loadFirstIPOfRangeFromFile()
	pingCount := len(ips) * pingTime
	bar := pb.StartNew(pingCount)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var data = make([]CloudflareIPData, 0)

	fmt.Println("开始延迟测速(TCP)：")

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
	fmt.Println("开始下载测速：")
	for i := 0; i < downloadTestCount; i++ {
		_, speed := DownloadSpeedHandler(data[i].ip)
		data[i].downloadSpeed = speed
		bar.Add(1)
	}
	bar.Finish()
	ExportCsv("./result.csv", data)
}
