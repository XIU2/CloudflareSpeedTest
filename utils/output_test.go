package utils

import (
	"encoding/json"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func sampleData() []CloudflareIPData {
	return []CloudflareIPData{
		{
			PingData:      &PingData{IP: &net.IPAddr{IP: net.ParseIP("104.16.1.1")}, Sended: 4, Received: 4, Delay: 45 * time.Millisecond, Colo: "FRA"},
			DownloadSpeed: 15.5 * 1024 * 1024,
		},
		{
			PingData:      &PingData{IP: &net.IPAddr{IP: net.ParseIP("104.16.1.2")}, Sended: 4, Received: 3, Delay: 120 * time.Millisecond, Colo: ""},
			DownloadSpeed: 0,
		},
	}
}

// 测试期间改写包级变量，结束后恢复，避免影响其他测试
func withOutput(t *testing.T, format, path string) {
	t.Helper()
	prevFormat, prevOutput := Format, Output
	Format, Output = format, path
	t.Cleanup(func() { Format, Output = prevFormat, prevOutput })
}

func TestExportJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "result.json")
	withOutput(t, FormatJSON, path)

	Export(sampleData())

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("读取 JSON 结果失败：%v", err)
	}
	var records []JSONRecord
	if err := json.Unmarshal(raw, &records); err != nil {
		t.Fatalf("JSON 格式不合法：%v\n%s", err, raw)
	}
	if len(records) != 2 {
		t.Fatalf("期望 2 条记录，实际 %d 条", len(records))
	}
	if records[0].IP != "104.16.1.1" || records[0].LossRate != 0 || records[0].LatencyMs != 45 || records[0].DownloadSpeedMB != 15.5 || records[0].Colo != "FRA" {
		t.Fatalf("第一条记录字段不符：%+v", records[0])
	}
	if records[1].Received != 3 || records[1].LossRate != 0.25 || records[1].Colo != "" {
		t.Fatalf("第二条记录字段不符：%+v", records[1])
	}
}

func TestExportIP(t *testing.T) {
	path := filepath.Join(t.TempDir(), "result.txt")
	withOutput(t, FormatIP, path)

	Export(sampleData())

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("读取 IP 结果失败：%v", err)
	}
	lines := strings.Split(strings.TrimSuffix(string(raw), "\n"), "\n")
	if len(lines) != 2 {
		t.Fatalf("期望 2 行，实际 %d 行：%q", len(lines), lines)
	}
	if lines[0] != "104.16.1.1" || lines[1] != "104.16.1.2" {
		t.Fatalf("IP 行内容不符：%q", lines)
	}
	if !strings.HasSuffix(string(raw), "\n") {
		t.Fatalf("结果应以换行结尾：%q", raw)
	}
}

func TestExportCsvUnchanged(t *testing.T) {
	path := filepath.Join(t.TempDir(), "result.csv")
	withOutput(t, FormatCsv, path)

	Export(sampleData())

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("读取 CSV 结果失败：%v", err)
	}
	want := "IP 地址,已发送,已接收,丢包率,平均延迟,下载速度(MB/s),地区码\n" +
		"104.16.1.1,4,4,0.00,45.00,15.50,FRA\n" +
		"104.16.1.2,4,3,0.25,120.00,0.00,N/A\n"
	if string(raw) != want {
		t.Fatalf("CSV 输出与原有格式不一致：\n得到 %q\n期望 %q", raw, want)
	}
}

func TestExportSkipsEmpty(t *testing.T) {
	for _, c := range []struct{ name, format, output string }{
		{"无结果", FormatJSON, filepath.Join(t.TempDir(), "empty.json")},
		{"空文件名", FormatJSON, ""},
		{"空格文件名", FormatIP, " "},
	} {
		t.Run(c.name, func(t *testing.T) {
			withOutput(t, c.format, c.output)
			Export(nil)
			if c.output == "" || c.output == " " {
				return
			}
			if _, err := os.Stat(c.output); !os.IsNotExist(err) {
				t.Fatalf("结果为空时不应创建文件：%s", c.output)
			}
		})
	}
}

func TestDefaultOutputName(t *testing.T) {
	for _, c := range []struct{ format, want string }{
		{FormatCsv, "result.csv"},
		{FormatJSON, "result.json"},
		{FormatIP, "result.txt"},
		{"xml", "result.csv"},
	} {
		prev := Format
		Format = c.format
		got := DefaultOutputName()
		Format = prev
		if got != c.want {
			t.Fatalf("格式 %s 的默认文件名应为 %s，实际 %s", c.format, c.want, got)
		}
	}
}

func TestValidFormat(t *testing.T) {
	for _, c := range []struct {
		format string
		want   bool
	}{
		{FormatCsv, true}, {FormatJSON, true}, {FormatIP, true}, {"xml", false}, {"", false},
	} {
		prev := Format
		Format = c.format
		got := ValidFormat()
		Format = prev
		if got != c.want {
			t.Fatalf("格式 %q 的校验结果应为 %v，实际 %v", c.format, c.want, got)
		}
	}
}
