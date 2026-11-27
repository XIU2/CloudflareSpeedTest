package utils

import (
	"bufio"
	"encoding/json"
	"log"
	"math"
	"os"
)

// 输出结果文件格式（由 -fmt 参数控制）
const (
	FormatCsv  = "csv"  // 兼容原有的 result.csv
	FormatJSON = "json" // JSON 数组，便于脚本解析
	FormatIP   = "ip"   // 每行一个 IP，便于直接用于代理配置等
)

// 非 csv 格式对应的默认输出文件名（未显式指定 -o 时使用）
const (
	defaultJsonOutput = "result.json"
	defaultIPOutput   = "result.txt"
)

// Format 是当前选择的输出结果文件格式
var Format = FormatCsv

// ValidFormat 判断 -fmt 的取值是否受支持
func ValidFormat() bool {
	return Format == FormatCsv || Format == FormatJSON || Format == FormatIP
}

// DefaultOutputName 返回当前格式对应的默认输出文件名
func DefaultOutputName() string {
	switch Format {
	case FormatJSON:
		return defaultJsonOutput
	case FormatIP:
		return defaultIPOutput
	}
	return defaultOutput
}

// JSONRecord 是 JSON 格式的单条测速结果，字段与 CSV 各列一一对应
type JSONRecord struct {
	IP              string  `json:"ip"`
	Sended          int     `json:"sent"`
	Received        int     `json:"received"`
	LossRate        float32 `json:"loss_rate"`
	LatencyMs       float64 `json:"latency_ms"`
	DownloadSpeedMB float64 `json:"download_speed_mbps"`
	Colo            string  `json:"colo"`
}

// Export 按 -fmt 指定的格式把测速结果写入文件
func Export(data []CloudflareIPData) {
	if noOutput() || len(data) == 0 {
		return
	}
	switch Format {
	case FormatJSON:
		exportJSON(data)
	case FormatIP:
		exportIP(data)
	default:
		ExportCsv(data)
	}
}

// 保留两位小数，与 CSV 的展示精度保持一致
func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

func exportJSON(data []CloudflareIPData) {
	records := make([]JSONRecord, 0, len(data))
	for _, v := range data {
		records = append(records, JSONRecord{
			IP:              v.IP.String(),
			Sended:          v.Sended,
			Received:        v.Received,
			LossRate:        float32(round2(float64(v.getLossRate()))),
			LatencyMs:       round2(v.Delay.Seconds() * 1000),
			DownloadSpeedMB: round2(v.DownloadSpeed / 1024 / 1024),
			Colo:            v.Colo,
		})
	}
	fp, err := os.Create(Output)
	if err != nil {
		log.Fatalf("创建文件[%s]失败：%v", Output, err)
		return
	}
	defer fp.Close()
	encoder := json.NewEncoder(fp)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(records); err != nil {
		log.Fatalf("写入文件[%s]失败：%v", Output, err)
	}
}

func exportIP(data []CloudflareIPData) {
	fp, err := os.Create(Output)
	if err != nil {
		log.Fatalf("创建文件[%s]失败：%v", Output, err)
		return
	}
	defer fp.Close()
	w := bufio.NewWriter(fp)
	for _, v := range data {
		if _, err := w.WriteString(v.IP.String() + "\n"); err != nil {
			log.Fatalf("写入文件[%s]失败：%v", Output, err)
		}
	}
	if err := w.Flush(); err != nil {
		log.Fatalf("写入文件[%s]失败：%v", Output, err)
	}
}
