package task

import (
	"net"
	"testing"
	"time"
)

func TestInitInterfaceEmpty(t *testing.T) {
	NetworkInterface = ""
	boundInterface = nil
	boundIPv4 = nil
	boundIPv6 = nil

	InitInterface()
	if boundInterface != nil || boundIPv4 != nil || boundIPv6 != nil {
		t.Fatalf("expected nil bound info when NetworkInterface is empty")
	}

	d := newDialer(nil, 5*time.Second)
	if d.LocalAddr != nil {
		t.Errorf("expected nil LocalAddr, got %v", d.LocalAddr)
	}
	if d.Control != nil {
		t.Errorf("expected nil Control, got non-nil")
	}
}

func TestInitInterfaceLoopback(t *testing.T) {
	ifaces, err := net.Interfaces()
	if err != nil || len(ifaces) == 0 {
		t.Skip("no network interfaces available")
	}

	// Pick first available interface for testing
	targetIface := ifaces[0].Name
	NetworkInterface = targetIface
	boundInterface = nil
	boundIPv4 = nil
	boundIPv6 = nil

	InitInterface()
	if boundInterface == nil {
		t.Fatalf("expected boundInterface to be set for %s", targetIface)
	}
	if boundInterface.Name != targetIface {
		t.Errorf("expected interface %s, got %s", targetIface, boundInterface.Name)
	}

	ipAddr := &net.IPAddr{IP: net.ParseIP("1.1.1.1")}
	d := newDialer(ipAddr, 3*time.Second)
	if d == nil {
		t.Fatalf("expected dialer to be created")
	}
}
