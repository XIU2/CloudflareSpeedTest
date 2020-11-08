package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/cheggaaa/pb/v3"
)

var version string
var disableDownload bool
var tcpPort int
var ipFile string
var outputFile string
var printResult int

func init() {
	var downloadSecond int64
	var printVersion bool
	const help = `
CloudflareSpeedTest
测试 Cloudflare CDN 所有 IP 的延迟和速度，获取最快 IP！
https://github.com/XIU2/CloudflareSpeedTest

参数：
    -n 500
        测速线程数量；数值越大速度越快，请勿超过 1000(结果误差大)；(默认 500)
    -t 4
        延迟测速次数；单个 IP 测速次数，为 1 时将过滤丢包的IP，TCP协议；(默认 4)
    -tp 443
        延迟测速端口；延迟测速 TCP 协议的端口；(默认 443)
    -dn 20
        下载测速数量；延迟测速并排序后，从最低延迟起下载测速数量，请勿太多(速度慢)；(默认 20)
    -dt 5
        下载测速时间；单个 IP 测速最长时间，单位：秒；(默认 5)
    -url https://cf.xiu2.xyz/Github/CloudflareSpeedTest.png
        下载测速地址；用来 Cloudflare CDN 测速的文件地址，如含有空格请加上引号；
    -p 20
        显示结果数量；测速后直接显示指定数量的结果，值为 0 时不显示结果直接退出；(默认 20)
    -f ip.txt
        IP 数据文件；如含有空格请加上引号；支持其他 CDN IP段，记得禁用下载测速；(默认 ip.txt)
    -o result.csv
        输出结果文件；如含有空格请加上引号；为空格时不输出结果文件(-o " ")；允许其他后缀；(默认 result.csv)
    -dd
        禁用下载测速；如果带上该参数就是禁用下载测速；(默认 启用)
    -v
        打印程序版本
    -h
        打印帮助说明

示例：
    注意：不需要加上所有参数（以下仅为示例），按需选择，参数前后顺序随意
    CloudflareST -n 500 -t 4 -dn 20 -dt 5
    CloudflareST -n 500 -t 4 -dn 20 -dt 5 -p 0 -f "ip.txt" -dd
    CloudflareST -n 500 -t 4 -dn 20 -dt 5 -p 20 -f "ip.txt" -o " " -dd
    CloudflareST -n 500 -t 4 -dn 20 -dt 5 -f "ip.txt" -o "result.csv" -dd
    CloudflareST -n 500 -t 4 -dn 20 -dt 5 -f "C:\abc\ip.txt" -o "C:\abc\result.csv" -dd
    CloudflareST -n 500 -t 4 -dn 20 -dt 5 -url https://cf.xiu2.xyz/Github/CloudflareSpeedTest.png`

	flag.IntVar(&pingRoutine, "n", 500, "测速线程数量")
	flag.IntVar(&pingTime, "t", 4, "延迟测速次数")
	flag.IntVar(&tcpPort, "tp", 443, "延迟测速端口")
	flag.IntVar(&downloadTestCount, "dn", 20, "下载测速数量")
	flag.Int64Var(&downloadSecond, "dt", 5, "下载测速时间")
	flag.StringVar(&url, "url", "https://cf.xiu2.xyz/Github/CloudflareSpeedTest.png", "下载测速地址")
	flag.IntVar(&printResult, "p", 20, "显示结果数量")
	flag.BoolVar(&disableDownload, "dd", false, "禁用下载测速")
	flag.StringVar(&ipFile, "f", "ip.txt", "IP 数据文件")
	flag.StringVar(&outputFile, "o", "result.csv", "输出结果文件")
	flag.BoolVar(&printVersion, "v", false, "打印程序版本")

	downloadTestTime = time.Duration(downloadSecond) * time.Second

	flag.Usage = func() { fmt.Print(help) }
	flag.Parse()
	if printVersion {
		println(version)
		os.Exit(0)
	}
	if pingRoutine <= 0 {
		pingRoutine = 500
	}
	if pingTime <= 0 {
		pingTime = 4
	}
	if tcpPort < 1 || tcpPort > 65535 {
		tcpPort = 443
	}
	if downloadTestCount <= 0 {
		downloadTestCount = 20
	}
	if downloadSecond <= 0 {
		downloadSecond = 10
	}
	if url == "" {
		url = "https://cf.xiu2.xyz/Github/CloudflareSpeedTest.png"
	}
	if printResult == 0 {
		printResult = 20
	}
	if ipFile == "" {
		ipFile = "ip.txt"
	}
	if outputFile == " " {
		outputFile = ""
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

	fmt.Println("开始延迟测速（模式：TCP，端口：" + strconv.Itoa(tcpPort) + "）：")
	control := make(chan bool, pingRoutine)
	for _, ip := range ips {
		wg.Add(1)
		control <- false
		handleProgress := handleProgressGenerator(bar) // 多线程进度条
		go tcpingGoroutine(&wg, &mu, ip, tcpPort, pingTime, &data, control, handleProgress)
	}
	wg.Wait()
	bar.Finish()

	sort.Sort(CloudflareIPDataSet(data)) // 排序

	// 下载测速
	if !disableDownload { // 如果禁用下载测速就跳过
		if len(data) > 0 { // IP数组长度(IP数量) 大于 0 时继续
			if len(data) < downloadTestCount { // 如果IP数组长度(IP数量) 小于 下载测速次数，则次数改为IP数
				downloadTestCount = len(data)
				fmt.Println("\n[信息] IP数量小于下载测速次数，下载测速次数改为IP数。\n")
			}
			bar = pb.Simple.Start(downloadTestCount)
			fmt.Println("开始下载测速：")
			for i := 0; i < downloadTestCount; i++ {
				_, speed := DownloadSpeedHandler(data[i].ip)
				data[i].downloadSpeed = speed
				bar.Add(1)
			}
			bar.Finish()
		} else {
			fmt.Println("\n[信息] IP数量为 0，跳过下载测速。")
		}
	}

	if outputFile != "" {
		ExportCsv(outputFile, data) // 输出结果到文件
	}

	// 直接输出结果
	if printResult > 0 { // 如果禁用下载测速就跳过
		dateString := convertToString(data) // 转为多维数组 [][]String
		if len(dateString) > 0 {            // IP数组长度(IP数量) 大于 0 时继续
			if len(dateString) < printResult { // 如果IP数组长度(IP数量) 小于  打印次数，则次数改为IP数量
				printResult = len(dateString)
				fmt.Println("\n[信息] IP数量小于显示结果数量，显示结果数量改为IP数量。\n")
			}
			fmt.Printf("%-16s%-5s%-5s%-5s%-6s%-11s\n", "IP 地址", "已发送", "已接收", "丢包率", "平均延迟", "下载速度 (MB/s)")
			for i := 0; i < printResult; i++ {
				fmt.Printf("%-18s%-8s%-8s%-8s%-10s%-15s\n", dateString[i][0], dateString[i][1], dateString[i][2], dateString[i][3], dateString[i][4], dateString[i][5])
			}
			if outputFile != "" {
				fmt.Printf("\n完整内容请查看 %v 文件。请按 回车键 或 Ctrl+C 退出。", outputFile)
			} else {
				fmt.Printf("\n请按 回车键 或 Ctrl+C 退出。")
			}
			var pause int
			fmt.Scanln(&pause)
		} else {
			fmt.Println("\n[信息] IP数量为 0，跳过输出结果。")
		}
	}
}
