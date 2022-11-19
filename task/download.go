package task

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"sort"
	"time"

	"CloudflareSpeedTest/utils"

	"github.com/VividCortex/ewma"
)

const (
	bufferSize                     = 1024
	defaultURL                     = "https://cf.xiu2.xyz/url"
	defaultTimeout                 = 10 * time.Second
	defaultDisableDownload         = false
	defaultTestNum                 = 10
	defaultMinSpeed        float64 = 0.0
)

var (
	// download test url
	URL = defaultURL
	// download timeout
	Timeout = defaultTimeout
	// disable download
	Disable = defaultDisableDownload

	TestCount = defaultTestNum
	MinSpeed  = defaultMinSpeed
)

func checkDownloadDefault() {
	if URL == "" {
		URL = defaultURL
	}
	if Timeout <= 0 {
		Timeout = defaultTimeout
	}
	if TestCount <= 0 {
		TestCount = defaultTestNum
	}
	if MinSpeed <= 0.0 {
		MinSpeed = defaultMinSpeed
	}
}

func TestDownloadSpeed(ipSet utils.PingDelaySet) (speedSet utils.DownloadSpeedSet) {
	checkDownloadDefault()
	if Disable {
		return utils.DownloadSpeedSet(ipSet)
	}
	if len(ipSet) <= 0 { // IP数组长度(IP数量) 大于 0 时才会继续下载测速
		fmt.Println("\n[信息] 延迟测速结果 IP 数量为 0，跳过下载测速。")
		return
	}
	testNum := TestCount
	if len(ipSet) < TestCount || MinSpeed > 0 { // 如果IP数组长度(IP数量) 小于下载测速数量（-dn），则次数修正为IP数
		testNum = len(ipSet)
	}
	if testNum < TestCount {
		TestCount = testNum
	}

	fmt.Printf("开始下载测速（下载速度下限：%.2f MB/s，下载测速数量：%d，下载测速队列：%d）：\n", MinSpeed, TestCount, testNum)
	bar := utils.NewBar(TestCount, "", "")
	for i := 0; i < testNum; i++ {
		speed := downloadHandler(ipSet[i].IP)
		ipSet[i].DownloadSpeed = speed
		// 在每个 IP 下载测速后，以 [下载速度下限] 条件过滤结果
		if speed >= MinSpeed*1024*1024 {
			bar.Grow(1, "")
			speedSet = append(speedSet, ipSet[i]) // 高于下载速度下限时，添加到新数组中
			if len(speedSet) == TestCount {       // 凑够满足条件的 IP 时（下载测速数量 -dn），就跳出循环
				break
			}
		}
	}
	bar.Done()
	if len(speedSet) == 0 { // 没有符合速度限制的数据，返回所有测试数据
		speedSet = utils.DownloadSpeedSet(ipSet)
	}
	// 按速度排序
	sort.Sort(speedSet)
	return
}

func getDialContext(ip *net.IPAddr) func(ctx context.Context, network, address string) (net.Conn, error) {
	var fakeSourceAddr string
	if isIPv4(ip.String()) {
		fakeSourceAddr = fmt.Sprintf("%s:%d", ip.String(), TCPPort)
	} else {
		fakeSourceAddr = fmt.Sprintf("[%s]:%d", ip.String(), TCPPort)
	}
	return func(ctx context.Context, network, address string) (net.Conn, error) {
		return (&net.Dialer{}).DialContext(ctx, network, fakeSourceAddr)
	}
}

// return download Speed
func downloadHandler(ip *net.IPAddr) float64 {
	client := &http.Client{
		Transport: &http.Transport{DialContext: getDialContext(ip)},
		Timeout:   Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) > 10 { // 限制最多重定向 10 次
				return http.ErrUseLastResponse
			}
			if req.Header.Get("Referer") == defaultURL { // 当使用默认下载测速地址时，重定向不携带 Referer
				req.Header.Del("Referer")
			}
			return nil
		},
	}
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return 0.0
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.80 Safari/537.36")

	response, err := client.Do(req)
	if err != nil {
		return 0.0
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return 0.0
	}
	timeStart := time.Now()           // 开始时间（当前）
	timeEnd := timeStart.Add(Timeout) // 加上下载测速时间得到的结束时间

	contentLength := response.ContentLength // 文件大小
	buffer := make([]byte, bufferSize)

	var (
		contentRead     int64 = 0
		timeSlice             = Timeout / 100
		timeCounter           = 1
		lastContentRead int64 = 0
	)

	var nextTime = timeStart.Add(timeSlice * time.Duration(timeCounter))
	e := ewma.NewMovingAverage()

	// 循环计算，如果文件下载完了（两者相等），则退出循环（终止测速）
	for contentLength != contentRead {
		currentTime := time.Now()
		if currentTime.After(nextTime) {
			timeCounter++
			nextTime = timeStart.Add(timeSlice * time.Duration(timeCounter))
			e.Add(float64(contentRead - lastContentRead))
			lastContentRead = contentRead
		}
		// 如果超出下载测速时间，则退出循环（终止测速）
		if currentTime.After(timeEnd) {
			break
		}
		bufferRead, err := response.Body.Read(buffer)
		if err != nil {
			if err != io.EOF { // 文件下载完了，或因网络等问题导致链接中断，则退出循环（终止测速）
				break
			}
			e.Add(float64(contentRead-lastContentRead) / (float64(nextTime.Sub(currentTime)) / float64(timeSlice)))
		}
		contentRead += int64(bufferRead)
	}
	return e.Value() / (Timeout.Seconds() / 120)
}
