package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/XIU2/CloudflareSpeedTest/task"
	"github.com/XIU2/CloudflareSpeedTest/utils"
)

var (
	version, versionNew string
)

func init() {
	var printVersion bool
	var help = `
CloudflareSpeedTest ` + version + `
测试 Cloudflare CDN 所有 IP 的延迟和速度，获取最快 IP (IPv4+IPv6)！
https://github.com/XIU2/CloudflareSpeedTest

参数：
    -n 200
        延迟测速线程；越多延迟测速越快，性能弱的设备 (如路由器) 请勿太高；(默认 200 最多 1000)
    -t 4
        延迟测速次数；单个 IP 延迟测速的次数；(默认 4 次)
    -dn 10
        下载测速数量；延迟测速并排序后，从最低延迟起下载测速的数量；(默认 10 个)
    -dt 10
        下载测速时间；单个 IP 下载测速最长时间，不能太短；(默认 10 秒)
    -tp 443
        指定测速端口；延迟测速/下载测速时使用的端口；(默认 443 端口)
    -url https://cf.xiu2.xyz/url
        指定测速地址；延迟测速(HTTPing)/下载测速时使用的地址，默认地址不保证可用性，建议自建；

    -httping
        切换测速模式；延迟测速模式改为 HTTP 协议，所用测试地址为 [-url] 参数；(默认 TCPing)
    -httping-code 200
        有效状态代码；HTTPing 延迟测速时网页返回的有效 HTTP 状态码，仅限一个；(默认 200 301 302)
    -cfcolo HKG,KHH,NRT,LAX,SEA,SJC,FRA,MAD
        匹配指定地区；地区名为当地机场三字码，英文逗号分隔，仅 HTTPing 模式可用；(默认 所有地区)

    -tl 200
        平均延迟上限；只输出低于指定平均延迟的 IP，各上下限条件可搭配使用；(默认 9999 ms)
    -tll 40
        平均延迟下限；只输出高于指定平均延迟的 IP；(默认 0 ms)
    -tlr 0.2
        丢包几率上限；只输出低于/等于指定丢包率的 IP，范围 0.00~1.00，0 过滤掉任何丢包的 IP；(默认 1.00)
    -sl 5
        下载速度下限；只输出高于指定下载速度的 IP，凑够指定数量 [-dn] 才会停止测速；(默认 0.00 MB/s)

    -p 10
        显示结果数量；测速后直接显示指定数量的结果，为 0 时不显示结果直接退出；(默认 10 个)
    -f ip.txt
        IP段数据文件；如路径含有空格请加上引号；支持其他 CDN IP段；(默认 ip.txt)
    -ip 1.1.1.1,2.2.2.2/24,2606:4700::/32
        指定IP段数据；直接通过参数指定要测速的 IP 段数据，英文逗号分隔；(默认 空)
    -o result.csv
        写入结果文件；如路径含有空格请加上引号；值为空时不写入文件 [-o ""]；(默认 result.csv)

    -dd
        禁用下载测速；禁用后测速结果会按延迟排序 (默认按下载速度排序)；(默认 启用)
    -allip
        测速全部的IP；对 IP 段中的每个 IP (仅支持 IPv4) 进行测速；(默认 每个 /24 段随机测速一个 IP)

    -v
        打印程序版本 + 检查版本更新
    -h
        打印帮助说明
`
	var minDelay, maxDelay, downloadTime int
	var maxLossRate float64
	flag.IntVar(&task.Routines, "n", 200, "延迟测速线程")
	flag.IntVar(&task.PingTimes, "t", 4, "延迟测速次数")
	flag.IntVar(&task.TestCount, "dn", 10, "下载测速数量")
	flag.IntVar(&downloadTime, "dt", 10, "下载测速时间")
	flag.IntVar(&task.TCPPort, "tp", 443, "指定测速端口")
	flag.StringVar(&task.URL, "url", "https://cf.xiu2.xyz/url", "指定测速地址")

	flag.BoolVar(&task.Httping, "httping", false, "切换测速模式")
	flag.IntVar(&task.HttpingStatusCode, "httping-code", 0, "有效状态代码")
	flag.StringVar(&task.HttpingCFColo, "cfcolo", "", "匹配指定地区")

	flag.IntVar(&maxDelay, "tl", 9999, "平均延迟上限")
	flag.IntVar(&minDelay, "tll", 0, "平均延迟下限")
	flag.Float64Var(&maxLossRate, "tlr", 1, "丢包几率上限")
	flag.Float64Var(&task.MinSpeed, "sl", 0, "下载速度下限")

	flag.IntVar(&utils.PrintNum, "p", 10, "显示结果数量")
	flag.StringVar(&task.IPFile, "f", "ip.txt", "IP段数据文件")
	flag.StringVar(&task.IPText, "ip", "", "指定IP段数据")
	flag.StringVar(&utils.Output, "o", "result.csv", "输出结果文件")

	flag.BoolVar(&task.Disable, "dd", false, "禁用下载测速")
	flag.BoolVar(&task.TestAll, "allip", false, "测速全部 IP")

	flag.BoolVar(&printVersion, "v", false, "打印程序版本")
	flag.Usage = func() { fmt.Print(help) }
	flag.Parse()

	if task.MinSpeed > 0 && time.Duration(maxDelay)*time.Millisecond == utils.InputMaxDelay {
		fmt.Println("[小提示] 在使用 [-sl] 参数时，建议搭配 [-tl] 参数，以避免因凑不够 [-dn] 数量而一直测速...")
	}
	utils.InputMaxDelay = time.Duration(maxDelay) * time.Millisecond
	utils.InputMinDelay = time.Duration(minDelay) * time.Millisecond
	utils.InputMaxLossRate = float32(maxLossRate)
	task.Timeout = time.Duration(downloadTime) * time.Second
	task.HttpingCFColomap = task.MapColoMap()

	if printVersion {
		println(version)
		fmt.Println("检查版本更新中...")
		checkUpdate()
		if versionNew != "" {
			fmt.Printf("*** 发现新版本 [%s]！请前往 [https://github.com/XIU2/CloudflareSpeedTest] 更新！ ***", versionNew)
		} else {
			fmt.Println("当前为最新版本 [" + version + "]！")
		}
		os.Exit(0)
	}
}

func main() {
	task.InitRandSeed() // 置随机数种子

	fmt.Printf("# XIU2/CloudflareSpeedTest %s \n\n", version)

	// 开始延迟测速 + 过滤延迟/丢包
	pingData := task.NewPing().Run().FilterDelay().FilterLossRate()
	// 开始下载测速
	speedData := task.TestDownloadSpeed(pingData)
	utils.ExportCsv(speedData) // 输出文件
	speedData.Print()          // 打印结果

	if versionNew != "" {
		fmt.Printf("\n*** 发现新版本 [%s]！请前往 [https://github.com/XIU2/CloudflareSpeedTest] 更新！ ***\n", versionNew)
	}
	endPrint()
}

func endPrint() {
	if utils.NoPrintResult() {
		return
	}
	if runtime.GOOS == "windows" { // 如果是 Windows 系统，则需要按下 回车键 或 Ctrl+C 退出（避免通过双击运行时，测速完毕后直接关闭）
		fmt.Printf("按下 回车键 或 Ctrl+C 退出。")
		fmt.Scanln()
	}
}

// 检查更新
func checkUpdate() {
	timeout := 10 * time.Second
	client := http.Client{Timeout: timeout}
	res, err := client.Get("https://api.xiu2.xyz/ver/cloudflarespeedtest.txt")
	if err != nil {
		return
	}
	// 读取资源数据 body: []byte
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	// 关闭资源流
	defer res.Body.Close()
	if string(body) != version {
		versionNew = string(body)
	}
}
