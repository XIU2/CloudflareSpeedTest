package task

import (
	"log"
	"net"
	"syscall"
)

var FWMark int

func newDialer() *net.Dialer {
	dialer := &net.Dialer{}
	if FWMark != 0 {
		dialer.Control = dialerController
	}
	return dialer
}

func dialerController(network, address string, c syscall.RawConn) error {
	return c.Control(func(fd uintptr) {
		err := syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_MARK, FWMark)
		if err != nil {
			log.Fatalln("failed to set fwmark", err)
		}
	})
}
