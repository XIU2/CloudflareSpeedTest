package utils

import (
	"net"
	"sort"
	"testing"
	"time"
)

func TestFilterDelayChecksEveryLossRateGroup(t *testing.T) {
	originalMaxDelay, originalMinDelay := InputMaxDelay, InputMinDelay
	t.Cleanup(func() {
		InputMaxDelay, InputMinDelay = originalMaxDelay, originalMinDelay
	})
	InputMaxDelay = 200 * time.Millisecond
	InputMinDelay = 0

	data := PingDelaySet{
		newPingDelay("192.0.2.1", 4, 50*time.Millisecond),
		newPingDelay("192.0.2.2", 4, 100*time.Millisecond),
		newPingDelay("192.0.2.3", 3, 80*time.Millisecond),
		newPingDelay("192.0.2.4", 3, 250*time.Millisecond),
		newPingDelay("192.0.2.5", 2, 70*time.Millisecond),
	}
	sort.Sort(data)

	got := data.FilterDelay()
	want := map[string]bool{
		"192.0.2.1": true,
		"192.0.2.2": true,
		"192.0.2.3": true,
		"192.0.2.5": true,
	}
	if len(got) != len(want) {
		t.Fatalf("FilterDelay() returned %d entries, want %d", len(got), len(want))
	}
	for _, item := range got {
		if !want[item.IP.String()] {
			t.Errorf("FilterDelay() unexpectedly returned %s", item.IP)
		}
		delete(want, item.IP.String())
	}
	for ip := range want {
		t.Errorf("FilterDelay() omitted eligible IP %s", ip)
	}
}

func newPingDelay(ip string, received int, delay time.Duration) CloudflareIPData {
	return CloudflareIPData{PingData: &PingData{
		IP:       &net.IPAddr{IP: net.ParseIP(ip)},
		Sended:   4,
		Received: received,
		Delay:    delay,
	}}
}
