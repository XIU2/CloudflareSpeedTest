package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/cheggaaa/pb/v3"
)

var version string
var disableDownload bool
var ipv6Mode bool
var tcpPort int
var ipFile string
var outputFile string
var printResultNum int
var timeLimit int
var speedLimit int

func init() {
	var downloadSecond int64
	var printVersion bool
	var help = `
CloudflareSpeedTest ` + version + `
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
    -tl 200
        延迟时间上限；只输出指定延迟时间以下的结果，数量为 -dn 参数的值，单位：ms；
    -sl 5
        下载速度下限；只输出指定下载速度以上的结果，数量为 -dn 参数的值，单位：MB/s；
    -p 20
        显示结果数量；测速后直接显示指定数量的结果，值为 0 时不显示结果直接退出；(默认 20)
    -f ip.txt
        IP 数据文件；如含有空格请加上引号；支持其他 CDN IP段，记得禁用下载测速；(默认 ip.txt)
    -o result.csv
        输出结果文件；如含有空格请加上引号；为空格时不输出结果文件(-o " ")；允许其他后缀；(默认 result.csv)
    -dd
        禁用下载测速；如果带上该参数就是禁用下载测速；(默认 启用下载测速)
    -ipv6
        IPv6 测速模式；请确保 IP 数据文件内只包含 IPv6 IP段，软件不支持同时测速 IPv4+IPv6；(默认 IPv4)
    -v
        打印程序版本
    -h
        打印帮助说明
`

	flag.IntVar(&pingRoutine, "n", 500, "测速线程数量")
	flag.IntVar(&pingTime, "t", 4, "延迟测速次数")
	flag.IntVar(&tcpPort, "tp", 443, "延迟测速端口")
	flag.IntVar(&downloadTestCount, "dn", 20, "下载测速数量")
	flag.Int64Var(&downloadSecond, "dt", 5, "下载测速时间")
	flag.StringVar(&url, "url", "https://cf.xiu2.xyz/Github/CloudflareSpeedTest.png", "下载测速地址")
	flag.IntVar(&timeLimit, "tl", 0, "延迟时间上限")
	flag.IntVar(&speedLimit, "sl", 0, "下载速度下限")
	flag.IntVar(&printResultNum, "p", 20, "显示结果数量")
	flag.BoolVar(&disableDownload, "dd", false, "禁用下载测速")
	flag.BoolVar(&ipv6Mode, "ipv6", false, "禁用下载测速")
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
	if timeLimit <= 0 {
		timeLimit = 9999
	}
	if speedLimit < 0 {
		speedLimit = 0
	}
	if printResultNum < 0 {
		printResultNum = 20
	}
	if ipFile == "" {
		ipFile = "ip.txt"
	}
	if outputFile == " " {
		outputFile = ""
	}
}

func main() {
	initRandSeed()                            // 置随机数种子
	failTime = pingTime                       // 设置接收次数
	ips := loadFirstIPOfRangeFromFile(ipFile) // 读入IP
	pingCount := len(ips) * pingTime          // 计算进度条总数（IP*测试次数）
	bar := pb.Simple.Start(pingCount)         // 进度条总数
	var wg sync.WaitGroup
	var mu sync.Mutex
	var data = make([]CloudflareIPData, 0)
	var data_2 = make([]CloudflareIPData, 0)

	fmt.Println("# XIU2/CloudflareSpeedTest " + version + "\n")
	if ipv6Mode {
		fmt.Println("开始延迟测速（模式：TCP IPv6，端口：" + strconv.Itoa(tcpPort) + "）：")
	} else {
		fmt.Println("开始延迟测速（模式：TCP IPv4，端口：" + strconv.Itoa(tcpPort) + "）：")
	}
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
				//fmt.Println("\n[信息] IP 数量小于下载测速次数（" + strconv.Itoa(downloadTestCount) + " < " + strconv.Itoa(len(data)) + "），下载测速次数改为IP数。\n")
				downloadTestCount = len(data)
			}
			var downloadTestCount_2 int // 临时的下载测速次数
			if timeLimit == 9999 && speedLimit == 0 {
				downloadTestCount_2 = downloadTestCount // 如果没有指定条件，则临时的下载次数变量为下载测速次数
				fmt.Println("开始下载测速：")
			} else if timeLimit > 0 || speedLimit >= 0 {
				downloadTestCount_2 = len(data) // 如果指定了任意一个条件，则临时的下载次数变量改为总数量
				fmt.Println("开始下载测速（延迟时间上限：" + strconv.Itoa(timeLimit) + " ms，下载速度下限：" + strconv.Itoa(speedLimit) + " MB/s）：")
			}
			bar = pb.Simple.Start(downloadTestCount_2)
			for i := 0; i < downloadTestCount_2; i++ {
				_, speed := DownloadSpeedHandler(data[i].ip)
				data[i].downloadSpeed = speed
				bar.Add(1)
				if int(data[i].pingTime) <= timeLimit && int(float64(speed)/1024/1024) >= speedLimit {
					data_2 = append(data_2, data[i])      // 延迟和速度均满足条件时，添加到新数组中
					if len(data_2) == downloadTestCount { // 满足条件的 IP =下载测速次数，则跳出循环
						break
					}
				} else if int(data[i].pingTime) > timeLimit {
					break
				}
			}
			bar.Finish()
		} else {
			fmt.Println("\n[信息] IP数量为 0，跳过下载测速。")
		}
	}

	if len(data_2) > 0 { // 如果该数字有内容，说明进行过指定条件的下载测速
		if outputFile != "" {
			ExportCsv(outputFile, data_2) // 输出结果到文件（指定延迟时间或下载速度的）
		}
		printResult(data_2) // 显示最快结果（指定延迟时间或下载速度的）
	} else {
		if outputFile != "" {
			ExportCsv(outputFile, data) // 输出结果到文件
		}
		printResult(data) // 显示最快结果
	}
}

// 显示最快结果
func printResult(data []CloudflareIPData) {
	sysType := runtime.GOOS
	if printResultNum > 0 { // 如果禁止直接输出结果就跳过
		dateString := convertToString(data) // 转为多维数组 [][]String
		if len(dateString) > 0 {            // IP数组长度(IP数量) 大于 0 时继续
			if len(dateString) < printResultNum { // 如果IP数组长度(IP数量) 小于  打印次数，则次数改为IP数量
				//fmt.Println("\n[信息] IP 数量小于显示结果数量（" + strconv.Itoa(printResultNum) + " < " + strconv.Itoa(len(dateString)) + "），显示结果数量改为IP数量。\n")
				printResultNum = len(dateString)
			}
			if ipv6Mode { // IPv6 太长了，所以需要调整一下间隔
				fmt.Printf("%-40s%-5s%-5s%-5s%-6s%-11s\n", "IP 地址", "已发送", "已接收", "丢包率", "平均延迟", "下载速度 (MB/s)")
				for i := 0; i < printResultNum; i++ {
					fmt.Printf("%-42s%-8s%-8s%-8s%-10s%-15s\n", ipPadding(dateString[i][0]), dateString[i][1], dateString[i][2], dateString[i][3], dateString[i][4], dateString[i][5])
				}
			} else {
				fmt.Printf("%-16s%-5s%-5s%-5s%-6s%-11s\n", "IP 地址", "已发送", "已接收", "丢包率", "平均延迟", "下载速度 (MB/s)")
				for i := 0; i < printResultNum; i++ {
					fmt.Printf("%-18s%-8s%-8s%-8s%-10s%-15s\n", ipPadding(dateString[i][0]), dateString[i][1], dateString[i][2], dateString[i][3], dateString[i][4], dateString[i][5])
				}
			}

			if sysType == "windows" { // 如果是 Windows 系统，则需要按下 回车键 或 Ctrl+C 退出
				if outputFile != "" {
					fmt.Printf("\n完整测速结果已写入 %v 文件，请使用记事本/表格软件查看。\n按下 回车键 或 Ctrl+C 退出。", outputFile)
				} else {
					fmt.Printf("\n按下 回车键 或 Ctrl+C 退出。")
				}
				var pause int
				fmt.Scanln(&pause)
			} else { // 其它系统直接退出
				if outputFile != "" {
					fmt.Println("\n完整测速结果已写入 " + outputFile + " 文件，请使用记事本/表格软件查看。")
				}
			}
		} else {
			fmt.Println("\n[信息] IP数量为 0，跳过输出结果。")
		}
	} else {
		fmt.Println("\n完整测速结果已写入 " + outputFile + " 文件，请使用记事本/表格软件查看。")
	}
}
