package cfdns

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
)

const testZoneID = "26afd8e5070162ce2b1fbfbe93934aa5"

type capturedRequest struct {
	Method      string
	Path        string
	Query       url.Values
	Auth        string
	ContentType string
	Body        string
}

// mockCF 用 httptest.Server 冒充 Cloudflare API，记录每个请求并按测试给定的逻辑应答
type mockCF struct {
	t        *testing.T
	server   *httptest.Server
	mu       sync.Mutex
	requests []capturedRequest
	respond  func(req capturedRequest) (int, interface{})
}

func newMockCF(t *testing.T, respond func(req capturedRequest) (int, interface{})) *mockCF {
	m := &mockCF{t: t, respond: respond}
	m.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("读取请求体失败: %v", err)
		}
		req := capturedRequest{
			Method:      r.Method,
			Path:        r.URL.Path,
			Query:       r.URL.Query(),
			Auth:        r.Header.Get("Authorization"),
			ContentType: r.Header.Get("Content-Type"),
			Body:        string(body),
		}
		m.mu.Lock()
		m.requests = append(m.requests, req)
		m.mu.Unlock()
		status, payload := m.respond(req)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(payload)
	}))
	t.Cleanup(m.server.Close)
	return m
}

// start 把 APIBase 指向打桩服务器（测试结束自动恢复），保证不会打到真实 Cloudflare API
func (m *mockCF) start() {
	old := APIBase
	APIBase = m.server.URL
	m.t.Cleanup(func() { APIBase = old })
}

func (m *mockCF) snapshot() []capturedRequest {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]capturedRequest(nil), m.requests...)
}

func (m *mockCF) countMethod(method string) int {
	n := 0
	for _, r := range m.snapshot() {
		if r.Method == method {
			n++
		}
	}
	return n
}

func unexpectedRequest(t *testing.T, req capturedRequest) (int, interface{}) {
	t.Errorf("意外的请求: %s %s", req.Method, req.Path)
	status, payload := cfError(http.StatusInternalServerError, "unexpected request in test")
	return status, payload
}

func cfSuccess(result interface{}) map[string]interface{} {
	return map[string]interface{}{"success": true, "errors": []interface{}{}, "result": result}
}

func cfError(status int, message string) (int, map[string]interface{}) {
	return status, map[string]interface{}{
		"success": false,
		"errors":  []map[string]string{{"code": "10000", "message": message}},
	}
}

func decodeBody(t *testing.T, raw string) dnsWriteBody {
	t.Helper()
	var body dnsWriteBody
	if err := json.Unmarshal([]byte(raw), &body); err != nil {
		t.Fatalf("请求体不是合法 JSON: %v (%s)", err, raw)
	}
	return body
}

// §4.1-5 zone 填根域名时能解析出 zone id（断言请求了 /zones?name=）
func TestUpdateResolvesZoneIDFromDomain(t *testing.T) {
	m := newMockCF(t, func(req capturedRequest) (int, interface{}) {
		switch {
		case req.Method == http.MethodGet && req.Path == "/zones":
			return http.StatusOK, cfSuccess([]map[string]string{{"id": testZoneID, "name": "example.com"}})
		case req.Method == http.MethodGet && strings.HasPrefix(req.Path, "/zones/"+testZoneID+"/dns_records"):
			return http.StatusOK, cfSuccess([]map[string]interface{}{})
		case req.Method == http.MethodPost && req.Path == "/zones/"+testZoneID+"/dns_records":
			return http.StatusOK, cfSuccess(map[string]string{"id": "new"})
		}
		return unexpectedRequest(t, req)
	})
	m.start()

	result, err := Update(Config{Zone: "example.com", Record: "cf.example.com", Token: "test-token", TTL: 60, Count: 1}, []string{"1.2.3.4"})
	if err != nil {
		t.Fatalf("Update 返回错误: %v", err)
	}
	requests := m.snapshot()
	if len(requests) < 2 {
		t.Fatalf("至少应有 2 次请求（查 zone + 查记录），实际: %d", len(requests))
	}
	if requests[0].Path != "/zones" {
		t.Errorf("第一个请求应查询 /zones，实际: %s", requests[0].Path)
	}
	if got := requests[0].Query.Get("name"); got != "example.com" {
		t.Errorf("/zones 请求应带 name=example.com，实际: %q", got)
	}
	if requests[1].Path != "/zones/"+testZoneID+"/dns_records" {
		t.Errorf("解析出的 Zone ID 应被用于后续请求路径，实际: %s", requests[1].Path)
	}
	if result.Created != 1 {
		t.Errorf("应创建 1 条记录，实际: %d", result.Created)
	}
}

// §4.1-5 既有记录为空 → 创建 N 条（断言 POST 次数 == N）
func TestUpdateCreatesMissingRecords(t *testing.T) {
	m := newMockCF(t, func(req capturedRequest) (int, interface{}) {
		switch {
		case req.Method == http.MethodGet && strings.HasPrefix(req.Path, "/zones/"+testZoneID+"/dns_records"):
			if got := req.Query.Get("type"); got != "A" {
				t.Errorf("查询记录类型应为 A，实际: %s", got)
			}
			return http.StatusOK, cfSuccess([]map[string]interface{}{})
		case req.Method == http.MethodPost && req.Path == "/zones/"+testZoneID+"/dns_records":
			return http.StatusOK, cfSuccess(map[string]string{"id": "new"})
		}
		return unexpectedRequest(t, req)
	})
	m.start()

	result, err := Update(Config{Zone: testZoneID, Record: "cf.example.com", Token: "test-token", TTL: 300, Count: 3}, []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"})
	if err != nil {
		t.Fatalf("Update 返回错误: %v", err)
	}
	if n := m.countMethod(http.MethodPost); n != 3 {
		t.Errorf("POST 次数应为 3，实际: %d", n)
	}
	if result.Created != 3 {
		t.Errorf("Result.Created 应为 3，实际: %d", result.Created)
	}
	if result.Skipped {
		t.Errorf("有写动作时 Skipped 不应为 true")
	}
	for _, r := range m.snapshot() {
		if r.Path == "/zones" {
			t.Errorf("Zone 已是 ID 形式时不应再请求 /zones")
		}
		if r.Query.Get("type") == "AAAA" {
			t.Errorf("只有 IPv4 时不应触碰 AAAA 记录")
		}
		if r.Method == http.MethodPost {
			if r.Auth != "Bearer test-token" {
				t.Errorf("Authorization 应为 Bearer test-token，实际: %q", r.Auth)
			}
			if r.ContentType != "application/json" {
				t.Errorf("Content-Type 应为 application/json，实际: %q", r.ContentType)
			}
		}
	}
	wantContents := "1.1.1.1,2.2.2.2,3.3.3.3"
	var contents []string
	for _, r := range m.snapshot() {
		if r.Method == http.MethodPost {
			body := decodeBody(t, r.Body)
			if body.Type != "A" || body.Name != "cf.example.com" || body.TTL != 300 || body.Proxied {
				t.Errorf("POST 请求体字段不符: %+v", body)
			}
			contents = append(contents, body.Content)
		}
	}
	if strings.Join(contents, ",") != wantContents {
		t.Errorf("POST 的 content 应按速度顺序为 %s，实际: %v", wantContents, contents)
	}
}

// §4.1-5 既有记录与目标完全相同 → 不发出任何写请求（幂等）
func TestUpdateSkipsWhenRecordsMatch(t *testing.T) {
	// 既有 contents 故意按与目标相反的顺序返回：按规格这是「集合相同」，应跳过
	existing := []map[string]interface{}{
		{"id": "r1", "type": "A", "name": "cf.example.com", "content": "2.2.2.2"},
		{"id": "r2", "type": "A", "name": "cf.example.com", "content": "1.1.1.1"},
	}
	m := newMockCF(t, func(req capturedRequest) (int, interface{}) {
		if req.Method == http.MethodGet && strings.HasPrefix(req.Path, "/zones/"+testZoneID+"/dns_records") {
			return http.StatusOK, cfSuccess(existing)
		}
		return unexpectedRequest(t, req)
	})
	m.start()

	result, err := Update(Config{Zone: testZoneID, Record: "cf.example.com", Token: "test-token", TTL: 60, Count: 2}, []string{"1.1.1.1", "2.2.2.2"})
	if err != nil {
		t.Fatalf("Update 返回错误: %v", err)
	}
	requests := m.snapshot()
	if len(requests) != 1 {
		t.Errorf("应只有 1 次 GET 查询请求，实际: %d", len(requests))
	}
	for _, method := range []string{http.MethodPut, http.MethodPost, http.MethodDelete} {
		if n := m.countMethod(method); n != 0 {
			t.Errorf("%s 请求次数应为 0，实际: %d", method, n)
		}
	}
	if !result.Skipped {
		t.Errorf("DNS 已是目标状态时应标记 Skipped=true")
	}
	if result.Created != 0 || result.Updated != 0 || result.Deleted != 0 {
		t.Errorf("跳过时不应有任何动作计数: %+v", result)
	}
}

// §4.1-5 既有 3 条、目标 1 条 → 1 次 PUT + 2 次 DELETE
func TestUpdatePutsAndDeletesExtras(t *testing.T) {
	existing := []map[string]interface{}{
		{"id": "r1", "type": "A", "name": "cf.example.com", "content": "9.9.9.9"},
		{"id": "r2", "type": "A", "name": "cf.example.com", "content": "8.8.8.8"},
		{"id": "r3", "type": "A", "name": "cf.example.com", "content": "7.7.7.7"},
	}
	var deletedIDs []string
	m := newMockCF(t, func(req capturedRequest) (int, interface{}) {
		switch {
		case req.Method == http.MethodGet && strings.HasPrefix(req.Path, "/zones/"+testZoneID+"/dns_records"):
			return http.StatusOK, cfSuccess(existing)
		case req.Method == http.MethodPut && strings.HasPrefix(req.Path, "/zones/"+testZoneID+"/dns_records/"):
			return http.StatusOK, cfSuccess(map[string]string{"id": "r1"})
		case req.Method == http.MethodDelete && strings.HasPrefix(req.Path, "/zones/"+testZoneID+"/dns_records/"):
			deletedIDs = append(deletedIDs, strings.TrimPrefix(req.Path, "/zones/"+testZoneID+"/dns_records/"))
			return http.StatusOK, cfSuccess(map[string]string{"id": "deleted"})
		}
		return unexpectedRequest(t, req)
	})
	m.start()

	result, err := Update(Config{Zone: testZoneID, Record: "cf.example.com", Token: "test-token", TTL: 300, Proxied: true, Count: 1}, []string{"1.1.1.1"})
	if err != nil {
		t.Fatalf("Update 返回错误: %v", err)
	}
	if n := m.countMethod(http.MethodPut); n != 1 {
		t.Errorf("PUT 次数应为 1，实际: %d", n)
	}
	if n := m.countMethod(http.MethodDelete); n != 2 {
		t.Errorf("DELETE 次数应为 2，实际: %d", n)
	}
	if n := m.countMethod(http.MethodPost); n != 0 {
		t.Errorf("POST 次数应为 0，实际: %d", n)
	}
	var putBody dnsWriteBody
	for _, r := range m.snapshot() {
		if r.Method == http.MethodPut {
			if r.Path != "/zones/"+testZoneID+"/dns_records/r1" {
				t.Errorf("PUT 应对齐第 1 条既有记录 r1，实际: %s", r.Path)
			}
			putBody = decodeBody(t, r.Body)
		}
	}
	if putBody.Content != "1.1.1.1" || putBody.Type != "A" || putBody.Name != "cf.example.com" || putBody.TTL != 300 || !putBody.Proxied {
		t.Errorf("PUT 请求体应携带 content/ttl/proxied: %+v", putBody)
	}
	if len(deletedIDs) != 2 || deletedIDs[0] != "r2" || deletedIDs[1] != "r3" {
		t.Errorf("应按顺序删除多余的 r2、r3，实际: %v", deletedIDs)
	}
	if result.Updated != 1 || result.Deleted != 2 || result.Skipped {
		t.Errorf("Result 计数不符: %+v", result)
	}
}

// §4.1-5 API 返回 403 → Update 返回 error，且 error 文本里含 CF 的 message
func TestUpdateAPIErrorPropagatesMessage(t *testing.T) {
	m := newMockCF(t, func(req capturedRequest) (int, interface{}) {
		return cfError(http.StatusForbidden, "Authentication error")
	})
	m.start()

	_, err := Update(Config{Zone: "example.com", Record: "cf.example.com", Token: "bad-token", TTL: 60, Count: 1}, []string{"1.2.3.4"})
	if err == nil {
		t.Fatal("API 返回 403 时应返回错误")
	}
	if err.Error() != "Cloudflare API 返回 403: Authentication error" {
		t.Errorf("错误信息应包含状态码与 CF 的 message，实际: %v", err)
	}
	if n := m.countMethod(http.MethodPut) + m.countMethod(http.MethodPost) + m.countMethod(http.MethodDelete); n != 0 {
		t.Errorf("解析 zone 失败后不应有写请求，实际: %d", n)
	}
}

// §4.1-5 ips 为空 → 返回 error 且零 API 请求
func TestUpdateEmptyIPs(t *testing.T) {
	m := newMockCF(t, func(req capturedRequest) (int, interface{}) {
		return unexpectedRequest(t, req)
	})
	m.start()

	result, err := Update(Config{Zone: testZoneID, Record: "cf.example.com", Token: "test-token", TTL: 60, Count: 1}, nil)
	if err == nil {
		t.Fatal("ips 为空时应返回错误")
	}
	if !strings.Contains(err.Error(), "没有可用 IP") {
		t.Errorf("错误信息应说明没有可用 IP，实际: %v", err)
	}
	if n := len(m.snapshot()); n != 0 {
		t.Errorf("ips 为空时应零 API 请求，实际: %d", n)
	}
	if result != (Result{}) {
		t.Errorf("应返回零值 Result，实际: %+v", result)
	}
}

// §4.1-5 IPv4/IPv6 混合 → 分别请求 type=A 与 type=AAAA
func TestUpdateMixedFamilies(t *testing.T) {
	m := newMockCF(t, func(req capturedRequest) (int, interface{}) {
		switch {
		case req.Method == http.MethodGet && strings.HasPrefix(req.Path, "/zones/"+testZoneID+"/dns_records"):
			return http.StatusOK, cfSuccess([]map[string]interface{}{})
		case req.Method == http.MethodPost && req.Path == "/zones/"+testZoneID+"/dns_records":
			return http.StatusOK, cfSuccess(map[string]string{"id": "new"})
		}
		return unexpectedRequest(t, req)
	})
	m.start()

	result, err := Update(Config{Zone: testZoneID, Record: "cf.example.com", Token: "test-token", TTL: 60, Count: 2},
		[]string{"1.1.1.1", "2606:4700::1", "2.2.2.2", "2606:4700::2"})
	if err != nil {
		t.Fatalf("Update 返回错误: %v", err)
	}
	var sawA, sawAAAA bool
	var aContents, aaaaContents []string
	for _, r := range m.snapshot() {
		if r.Method == http.MethodGet && strings.HasPrefix(r.Path, "/zones/"+testZoneID+"/dns_records") {
			switch r.Query.Get("type") {
			case "A":
				sawA = true
			case "AAAA":
				sawAAAA = true
			}
		}
		if r.Method == http.MethodPost {
			body := decodeBody(t, r.Body)
			switch body.Type {
			case "A":
				aContents = append(aContents, body.Content)
			case "AAAA":
				aaaaContents = append(aaaaContents, body.Content)
			}
		}
	}
	if !sawA || !sawAAAA {
		t.Errorf("混合 IP 应分别按 type=A 与 type=AAAA 查询（A: %v, AAAA: %v）", sawA, sawAAAA)
	}
	if strings.Join(aContents, ",") != "1.1.1.1,2.2.2.2" {
		t.Errorf("A 记录应写入 2 条 IPv4 且保持速度顺序，实际: %v", aContents)
	}
	if strings.Join(aaaaContents, ",") != "2606:4700::1,2606:4700::2" {
		t.Errorf("AAAA 记录应写入 2 条 IPv6 且保持速度顺序，实际: %v", aaaaContents)
	}
	if result.Created != 4 {
		t.Errorf("应创建 4 条记录，实际: %d", result.Created)
	}
}
