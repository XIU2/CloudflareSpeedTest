//go:build !linux

package task

import (
	"net"
)

var FWMark int

func newDialer() *net.Dialer {
	return &net.Dialer{}
}
