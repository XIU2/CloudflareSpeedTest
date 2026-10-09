// Package cfdns 把测速得到的最优 IP 同步到 Cloudflare 托管的 DNS（A/AAAA 记录）。
package cfdns

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"sort"
	"time"

	"github.com/XIU2/CloudflareSpeedTest/utils"
)

// APIBase 独立成包级变量是为了让单元测试能把请求打到本地 httptest.Server，而不是真实的 Cloudflare API
var APIBase = "https://api.cloudflare.com/client/v4"

// init 允许用环境变量覆盖 API 地址：测试时指向本地打桩服务，也可以指向兼容 Cloudflare API 的自建端点
func init() {
	if base := os.Getenv("CF_API_BASE"); base != "" {
		APIBase = base
	}
}

// httpClient 全局复用同一个连接池，避免每次更新重复建连
var httpClient = &http.Client{Timeout: 15 * time.Second}

// Config 一次 DNS 更新所需的全部参数
type Config struct {
	Zone     string // 根域名或 Zone ID
	Record   string // 完整记录名，如 cf.example.com
	Token    string
	TTL      int
	Proxied  bool
	Count    int           // 同名同类型最多写几条
	Resolver *net.Resolver // 为空时用系统默认解析器；仅用于 Cloudflare API 的域名解析
}

// clientFor 返回本次更新使用的 HTTP client；指定了 Resolver 时用它解析域名
func clientFor(cfg Config) *http.Client {
	if cfg.Resolver == nil {
		return httpClient
	}
	return &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{Timeout: 15 * time.Second, Resolver: cfg.Resolver}).DialContext,
		},
	}
}

// Result 一次更新实际发生的动作汇总
type Result struct {
	Created int
	Updated int
	Deleted int
	Skipped bool // true = DNS 已经是目标状态，什么都没改
}

type cfZone struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type dnsRecord struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
	Proxied bool   `json:"proxied"`
}

// dnsWriteBody PUT/POST 的请求体，字段与 Cloudflare API 的记录结构保持一致
type dnsWriteBody struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
	Proxied bool   `json:"proxied"`
}

type cfResponse struct {
	Success bool `json:"success"`
	Errors  []struct {
		Message string `json:"message"`
	} `json:"errors"`
	Result json.RawMessage `json:"result"`
}

// Update 把 ips 里最快的若干个 IP 同步到 DNS。ips 必须已按速度从高到低排好序。
// 返回实际写入/删除的动作数量与 error。
func Update(cfg Config, ips []string) (Result, error) {
	var result Result
	if len(ips) == 0 {
		return result, errors.New("DNS 更新跳过：没有可用 IP")
	}
	client := clientFor(cfg)
	zoneID, err := resolveZone(cfg, client)
	if err != nil {
		return result, err
	}
	// IPv4/IPv6 各自独立走一遍同步流程；只有一种家族时，另一种类型的既有记录保持不动
	ipv4, ipv6 := splitByFamily(ips)
	anySkipped := false
	if len(ipv4) > 0 {
		skipped, err := syncFamily(cfg, client, zoneID, "A", ipv4, &result)
		if err != nil {
			return result, err
		}
		anySkipped = anySkipped || skipped
	}
	if len(ipv6) > 0 {
		skipped, err := syncFamily(cfg, client, zoneID, "AAAA", ipv6, &result)
		if err != nil {
			return result, err
		}
		anySkipped = anySkipped || skipped
	}
	// 只要发生过任何动作就不是「跳过」；全部分家族都无动作且都命中幂等判断才算
	result.Skipped = anySkipped && result.Created == 0 && result.Updated == 0 && result.Deleted == 0
	return result, nil
}

// resolveZone 把根域名解析成 Zone ID；形如 32 位十六进制的值直接当 Zone ID 用，省一次 API 调用
func resolveZone(cfg Config, client *http.Client) (string, error) {
	if isZoneID(cfg.Zone) {
		return cfg.Zone, nil
	}
	var zones []cfZone
	if err := apiRequest(client, http.MethodGet, "/zones?name="+url.QueryEscape(cfg.Zone), cfg.Token, nil, &zones); err != nil {
		return "", err
	}
	if len(zones) == 0 {
		return "", fmt.Errorf("找不到域名 %s 对应的 Cloudflare Zone（请确认 token 有 Zone:Read 权限）", cfg.Zone)
	}
	return zones[0].ID, nil
}

// syncFamily 把一个 IP 家族按速度序取前 Count 条，对齐同步到同一种类型的记录上，
// 返回值表示该家族是否命中「已是最新」的幂等判断
func syncFamily(cfg Config, client *http.Client, zoneID, recordType string, familyIPs []string, result *Result) (bool, error) {
	count := cfg.Count
	if count < 1 {
		count = 1
	}
	targets := familyIPs
	if len(targets) > count {
		targets = targets[:count]
	}
	existing, err := listRecords(cfg, client, zoneID, recordType)
	if err != nil {
		return false, err
	}
	// 内容集合完全一致说明 DNS 已是目标状态；跳过才能保证 cron 反复跑不产生无谓的写请求
	if sameContents(existing, targets) {
		utils.Yellow.Println("[信息] DNS 记录已是最新，跳过更新。")
		return true, nil
	}
	// 第 i 条目标 IP 对齐第 i 条既有记录：有就 PUT（同时带上 ttl/proxied），没有就 POST 新建
	for i, ip := range targets {
		body := dnsWriteBody{Type: recordType, Name: cfg.Record, Content: ip, TTL: cfg.TTL, Proxied: cfg.Proxied}
		if i < len(existing) {
			if err := apiRequest(client, http.MethodPut, "/zones/"+zoneID+"/dns_records/"+existing[i].ID, cfg.Token, body, nil); err != nil {
				return false, err
			}
			result.Updated++
			utils.Yellow.Printf("[信息] 更新 %s 记录: %s\n", recordType, ip)
			continue
		}
		if err := apiRequest(client, http.MethodPost, "/zones/"+zoneID+"/dns_records", cfg.Token, body, nil); err != nil {
			return false, err
		}
		result.Created++
		utils.Green.Printf("[信息] 创建 %s 记录: %s\n", recordType, ip)
	}
	// 目标数量少于既有数量时，多余的记录删除
	if len(existing) > len(targets) {
		for _, record := range existing[len(targets):] {
			if err := apiRequest(client, http.MethodDelete, "/zones/"+zoneID+"/dns_records/"+record.ID, cfg.Token, nil, nil); err != nil {
				return false, err
			}
			result.Deleted++
			utils.Yellow.Printf("[信息] 删除 %s 记录: %s\n", recordType, record.Content)
		}
	}
	return false, nil
}

func listRecords(cfg Config, client *http.Client, zoneID, recordType string) ([]dnsRecord, error) {
	var records []dnsRecord
	listPath := fmt.Sprintf("/zones/%s/dns_records?type=%s&name=%s&per_page=100", zoneID, recordType, url.QueryEscape(cfg.Record))
	if err := apiRequest(client, http.MethodGet, listPath, cfg.Token, nil, &records); err != nil {
		return nil, err
	}
	return records, nil
}

// apiRequest 发起一次 Cloudflare API 请求；out 非 nil 时把 result 解析进去
func apiRequest(client *http.Client, method, apiPath, token string, body interface{}, out interface{}) error {
	var payload io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return err
		}
		payload = bytes.NewReader(encoded)
	}
	request, err := http.NewRequest(method, APIBase+apiPath, payload)
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("读取 Cloudflare API 响应失败: %w", err)
	}
	var parsed cfResponse
	// 非 JSON 响应（如网关返回的 HTML 错误页）也要按 HTTP 状态码报错，所以忽略解析失败
	_ = json.Unmarshal(data, &parsed)
	if response.StatusCode < 200 || response.StatusCode >= 300 || !parsed.Success {
		message := "未知错误"
		if len(parsed.Errors) > 0 && parsed.Errors[0].Message != "" {
			message = parsed.Errors[0].Message
		}
		return fmt.Errorf("Cloudflare API 返回 %d: %s", response.StatusCode, message)
	}
	if out != nil && len(parsed.Result) > 0 {
		if err := json.Unmarshal(parsed.Result, out); err != nil {
			return fmt.Errorf("Cloudflare API result 解析失败: %w", err)
		}
	}
	return nil
}

// splitByFamily 按 IP 家族分组并保持原有的速度顺序，无法解析的 IP 直接忽略
func splitByFamily(ips []string) (ipv4, ipv6 []string) {
	for _, s := range ips {
		ip := net.ParseIP(s)
		if ip == nil {
			continue
		}
		if ip.To4() != nil {
			ipv4 = append(ipv4, s)
		} else {
			ipv6 = append(ipv6, s)
		}
	}
	return
}

// sameContents 判断既有记录的 content 集合与目标集合是否完全一致（与顺序无关）
func sameContents(records []dnsRecord, targets []string) bool {
	if len(records) != len(targets) {
		return false
	}
	existing := make([]string, 0, len(records))
	for _, record := range records {
		existing = append(existing, record.Content)
	}
	sortedTargets := append([]string(nil), targets...)
	sort.Strings(existing)
	sort.Strings(sortedTargets)
	for i := range existing {
		if existing[i] != sortedTargets[i] {
			return false
		}
	}
	return true
}

// isZoneID 形如 32 位十六进制（不分大小写）的 zone 值按 Zone ID 直接使用
func isZoneID(s string) bool {
	if len(s) != 32 {
		return false
	}
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}
