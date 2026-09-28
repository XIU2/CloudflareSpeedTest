//go:build !darwin && !linux

package task

import (
	"net"
	"syscall"
)

// bindInterfaceControl is a no-op fallback on platforms where interface socket options are not supported directly
func bindInterfaceControl(iface *net.Interface, isIPv4 bool) func(network, address string, c syscall.RawConn) error {
	return nil
}
