package utils

import (
	"net"
	"sort"
	"testing"
	"time"
)

func mustIP(s string) *net.IPAddr {
	return &net.IPAddr{IP: net.ParseIP(s)}
}

// 复现 -tl 过滤误杀问题：
// Ping.Run() 返回的数据是按「丢包率优先、同丢包率内按延迟升序」排序的（见 PingDelaySet.Less），
// 因此 FilterDelay() 中遇到第一个延迟超上限的 IP 就 break 是错误的——
// 它会连带丢弃排在后面其他丢包率分组中的低延迟 IP。
func TestFilterDelay_NotAbortAtMaxDelay(t *testing.T) {
	oldMax, oldMin := InputMaxDelay, InputMinDelay
	defer func() { InputMaxDelay, InputMinDelay = oldMax, oldMin }()

	InputMaxDelay = 200 * time.Millisecond
	InputMinDelay = 0

	// 与 Ping.Run() 一致的排序（sort.Sort(PingDelaySet)，丢包率优先）后应为：
	// [0%丢包 100ms, 0%丢包 500ms, 25%丢包 50ms]
	// 其中 50ms 的 IP 延迟满足 -tl 200 条件，但因排在 500ms（触发 break 的 IP）
	// 之后而被错误丢弃——这正是 bug 所在
	s := PingDelaySet{
		{PingData: &PingData{IP: mustIP("1.1.1.1"), Sended: 4, Received: 4, Delay: 500 * time.Millisecond}}, // 0% 丢包
		{PingData: &PingData{IP: mustIP("2.2.2.2"), Sended: 4, Received: 4, Delay: 100 * time.Millisecond}}, // 0% 丢包
		{PingData: &PingData{IP: mustIP("3.3.3.3"), Sended: 4, Received: 3, Delay: 50 * time.Millisecond}},  // 25% 丢包
	}
	sort.Sort(s)

	got := s.FilterDelay()
	if len(got) != 2 {
		t.Fatalf("FilterDelay() 保留了 %d 个 IP，期望 2 个：100ms 无丢包，以及延迟仅 50ms 满足条件的低延迟 IP", len(got))
	}
	for _, v := range got {
		if v.Delay > InputMaxDelay || v.Delay < InputMinDelay {
			t.Fatalf("IP %s 平均延迟 %v 超出 -tl %v / -tll %v 条件却被保留", v.IP, v.Delay, InputMinDelay, InputMaxDelay)
		}
	}
	if got[0].Delay != 100*time.Millisecond || got[1].Delay != 50*time.Millisecond {
		t.Fatalf("过滤结果延迟不符，期望 100ms、50ms 两个 IP，实际为 %v 与 %v", got[0].Delay, got[1].Delay)
	}
}

// 仅指定 -tll（延迟下限）时，低于下限的 IP 跳过、其余正常保留
func TestFilterDelay_MinDelaySkip(t *testing.T) {
	oldMax, oldMin := InputMaxDelay, InputMinDelay
	defer func() { InputMaxDelay, InputMinDelay = oldMax, oldMin }()

	InputMaxDelay = 9999 * time.Millisecond
	InputMinDelay = 80 * time.Millisecond

	s := PingDelaySet{
		{PingData: &PingData{IP: mustIP("1.1.1.1"), Sended: 4, Received: 4, Delay: 50 * time.Millisecond}},
		{PingData: &PingData{IP: mustIP("2.2.2.2"), Sended: 4, Received: 4, Delay: 100 * time.Millisecond}},
		{PingData: &PingData{IP: mustIP("3.3.3.3"), Sended: 4, Received: 4, Delay: 120 * time.Millisecond}},
	}
	sort.Sort(s)

	got := s.FilterDelay()
	if len(got) != 2 {
		t.Fatalf("FilterDelay() 保留了 %d 个 IP，期望 2 个（100ms 与 120ms，50ms 应被 -tll 过滤）", len(got))
	}
}
