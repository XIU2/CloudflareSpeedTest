package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/cheggaaa/pb/v3"
)

var version, ipFile, outputFile, versionNew string
var disableDownload, ipv6Mode, allip bool
var tcpPort, printResultNum, timeLimit, speedLimit, downloadSecond int

func init() {
	var printVersion bool
	var help = `
CloudflareSpeedTest ` + version + `
测试 Cloudflare CDN 所有 IP 的延迟和速度，获取最快 IP (IPv4+IPv6)！
https://github.com/XIU2/CloudflareSpeedTest

参数：
    -n 500
        测速线程数量；越多测速越快，性能弱的设备 (如路由器) 请适当调低；(默认 500 最多 1000)
    -t 4
        延迟测速次数；单个 IP 延迟测速次数，为 1 时将过滤丢包的IP，TCP协议；(默认 4)
    -tp 443
        延迟测速端口；延迟测速 TCP 协议的端口；(默认 443)
    -dn 20
        下载测速数量；延迟测速并排序后，从最低延迟起下载测速的数量；(默认 20)
    -dt 10
        下载测速时间；单个 IP 下载测速最长时间，单位：秒；(默认 10)
    -url https://cf.xiu2.xyz/Github/CloudflareSpeedTest.png
        下载测速地址；用来下载测速的 Cloudflare CDN 文件地址，如地址含有空格请加上引号；
    -tl 200
        平均延迟上限；只输出低于指定平均延迟的 IP，与下载速度下限搭配使用；(默认 9999 ms)
    -sl 5
        下载速度下限；只输出高于指定下载速度的 IP，凑够指定数量 [-dn] 才会停止测速；(默认 0 MB/s)
    -p 20
        显示结果数量；测速后直接显示指定数量的结果，为 0 时不显示结果直接退出；(默认 20)
    -f ip.txt
        IP段数据文件；如路径含有空格请加上引号；支持其他 CDN IP段；(默认 ip.txt)
    -o result.csv
        输出结果文件；如路径含有空格请加上引号；值为空格时不输出 [-o " "]；(默认 result.csv)
    -dd
        禁用下载测速；禁用后测速结果会按延迟排序 (默认按下载速度排序)；(默认 启用)
    -ipv6
        IPv6测速模式；确保 IP 段数据文件内只包含 IPv6 IP段，软件不支持同时测速 IPv4+IPv6；(默认 IPv4)
    -allip
        测速全部的IP；对 IP 段中的每个 IP (仅支持 IPv4) 进行测速；(默认 每个 IP 段随机测速一个 IP)
    -v
        打印程序版本+检查版本更新
    -h
        打印帮助说明
`

	flag.IntVar(&pingRoutine, "n", 500, "测速线程数量")
	flag.IntVar(&pingTime, "t", 4, "延迟测速次数")
	flag.IntVar(&tcpPort, "tp", 443, "延迟测速端口")
	flag.IntVar(&downloadTestCount, "dn", 20, "下载测速数量")
	flag.IntVar(&downloadSecond, "dt", 10, "下载测速时间")
	flag.StringVar(&url, "url", "https://cf.xiu2.xyz/Github/CloudflareSpeedTest.png", "下载测速地址")
	flag.IntVar(&timeLimit, "tl", 9999, "延迟时间上限")
	flag.IntVar(&speedLimit, "sl", 0, "下载速度下限")
	flag.IntVar(&printResultNum, "p", 20, "显示结果数量")
	flag.BoolVar(&disableDownload, "dd", false, "禁用下载测速")
	flag.BoolVar(&ipv6Mode, "ipv6", false, "禁用下载测速")
	flag.BoolVar(&allip, "allip", false, "测速全部 IP")
	flag.StringVar(&ipFile, "f", "ip.txt", "IP 数据文件")
	flag.StringVar(&outputFile, "o", "result.csv", "输出结果文件")
	flag.BoolVar(&printVersion, "v", false, "打印程序版本")

	flag.Usage = func() { fmt.Print(help) }
	flag.Parse()
	if printVersion {
		println(version)
		fmt.Println("检查版本更新中...")
		checkUpdate()
		if versionNew != "" {
			fmt.Println("发现新版本 [" + versionNew + "]！请前往 [https://github.com/XIU2/CloudflareSpeedTest] 更新！")
		} else {
			fmt.Println("当前为最新版本 [" + version + "]！")
		}
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
	go checkUpdate()                          // 检查版本更新
	initRandSeed()                            // 置随机数种子
	failTime = pingTime                       // 设置接收次数
	ips := loadFirstIPOfRangeFromFile(ipFile) // 读入IP
	pingCount := len(ips) * pingTime          // 计算进度条总数（IP*测试次数）
	bar := pb.Simple.Start(pingCount)         // 进度条总数
	var wg sync.WaitGroup
	var mu sync.Mutex
	var data = make([]CloudflareIPData, 0)
	var data_2 = make([]CloudflareIPData, 0)
	downloadTestTime = time.Duration(downloadSecond) * time.Second

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
				downloadTestCount_2 = downloadTestCount // 如果没有指定条件，则临时变量为下载测速次数
				fmt.Println("开始下载测速：")
			} else if timeLimit > 0 || speedLimit >= 0 {
				downloadTestCount_2 = len(data) // 如果指定了任意一个条件，则临时变量改为总数量
				fmt.Println("开始下载测速（延迟时间上限：" + strconv.Itoa(timeLimit) + " ms，下载速度下限：" + strconv.Itoa(speedLimit) + " MB/s）：")
			}
			bar = pb.Simple.Start(downloadTestCount)
			for i := 0; i < downloadTestCount_2; i++ {
				_, speed := DownloadSpeedHandler(data[i].ip)
				data[i].downloadSpeed = speed
				if int(data[i].pingTime) <= timeLimit && int(float64(speed)/1024/1024) >= speedLimit {
					data_2 = append(data_2, data[i]) // 延迟和速度均满足条件时，添加到新数组中
					bar.Add(1)
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
		sort.Sort(CloudflareIPDataSetD(data_2)) // 排序
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

			if versionNew != "" {
				fmt.Println("\n发现新版本 [" + versionNew + "]！请前往 [https://github.com/XIU2/CloudflareSpeedTest] 更新！")
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

// 检查更新
func checkUpdate() {
	timeout := time.Duration(10 * time.Second)
	client := http.Client{Timeout: timeout}
	res, err := client.Get("https://api.xiuer.pw/ver/cloudflarespeedtest.txt")
	if err == nil {
		// 读取资源数据 body: []byte
		body, err := ioutil.ReadAll(res.Body)
		// 关闭资源流
		res.Body.Close()
		if err == nil {
			if string(body) != version {
				versionNew = string(body)
			}
		}
	}
}
