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
var disableDownload bool
var ipFile string

func init() {
	var downloadSecond int64
	var printVersion bool
	const help = `
CloudflareSpeedTest
测试 Cloudflare CDN 所有 IP 的延迟和速度，获取最快 IP！
https://github.com/XIU2/CloudflareSpeedTest

参数：
    -n 500
        测速线程数量；数值越大速度越快，请勿超过1000(结果误差大)；(默认 500)
    -t 4
        延迟测速次数；单个 IP 测速次数，为 1 时将过滤丢包的IP，TCP协议；(默认 4)
    -dn 20
        下载测速数量；延迟测速并排序后，从最低延迟起测试下载速度的数量，请勿太多(速度慢)；(默认 20)
    -dt 10
        下载测速时间；单个 IP 测速最长时间，单位：秒；(默认 10)
    -f ip.txt
        IP 数据文件；支持相对路径和绝对路径，如果包含空格请前后加上引号；(默认 ip.txt)
    -dd
        禁用下载测速；如果带上该参数就是禁用下载测速；(默认 启用)
    -v
        打印程序版本
    -h
        打印帮助说明

示例：
    CloudflareST.exe -n 500 -t 4 -dn 20 -dt 10
    CloudflareST.exe -n 500 -t 4 -dn 20 -dt 10 -f "C:\abc\ip.txt" -dd`

	flag.IntVar(&pingRoutine, "n", 500, "测速线程数量")
	flag.IntVar(&pingTime, "t", 4, "延迟测速次数")
	flag.IntVar(&downloadTestCount, "dn", 20, "下载测速数量")
	flag.Int64Var(&downloadSecond, "dt", 10, "下载测速时间")
	flag.BoolVar(&disableDownload, "dd", false, "禁用下载测速")
	flag.StringVar(&ipFile, "f", "ip.txt", "IP 数据文件")
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
	initipEndWith()                           // 随机数
	failTime = pingTime                       // 设置接收次数
	ips := loadFirstIPOfRangeFromFile(ipFile) // 读入IP
	pingCount := len(ips) * pingTime          // 计算进度条总数（IP*测试次数）
	bar := pb.Full.Start(pingCount)           // 进度条总数
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

	sort.Sort(CloudflareIPDataSet(data)) // 排序
	if !disableDownload {                // 如果禁用下载测速就跳过
		bar = pb.Simple.Start(downloadTestCount)
		fmt.Println("开始下载测速：")
		for i := 0; i < downloadTestCount; i++ {
			_, speed := DownloadSpeedHandler(data[i].ip)
			data[i].downloadSpeed = speed
			bar.Add(1)
		}
		bar.Finish()
	}
	ExportCsv("./result.csv", data) // 输出结果
}
