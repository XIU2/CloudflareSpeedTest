//go:build linux

package task

import (
	"net"
	"syscall"

	"golang.org/x/sys/unix"
)

// bindInterfaceControl returns a Control function to bind socket to Linux network device
func bindInterfaceControl(iface *net.Interface, isIPv4 bool) func(network, address string, c syscall.RawConn) error {
	if iface == nil {
		return nil
	}
	return func(network, address string, c syscall.RawConn) error {
		var operr error
		err := c.Control(func(fd uintptr) {
			operr = unix.BindToDevice(int(fd), iface.Name)
		})
		if err != nil {
			return err
		}
		return operr
	}
}
