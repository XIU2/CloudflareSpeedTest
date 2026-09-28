//go:build darwin

package task

import (
	"net"
	"syscall"

	"golang.org/x/sys/unix"
)

// bindInterfaceControl returns a Control function to bind socket to Darwin network interface
func bindInterfaceControl(iface *net.Interface, isIPv4 bool) func(network, address string, c syscall.RawConn) error {
	if iface == nil {
		return nil
	}
	return func(network, address string, c syscall.RawConn) error {
		var operr error
		err := c.Control(func(fd uintptr) {
			if isIPv4 {
				operr = unix.SetsockoptInt(int(fd), unix.IPPROTO_IP, unix.IP_BOUND_IF, iface.Index)
			} else {
				operr = unix.SetsockoptInt(int(fd), unix.IPPROTO_IPV6, unix.IPV6_BOUND_IF, iface.Index)
			}
		})
		if err != nil {
			return err
		}
		return operr
	}
}
