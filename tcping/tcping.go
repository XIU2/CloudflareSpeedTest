package tcping

import (
	"net"
	"sync"

	"CloudflareSpeedTest/utils"
)

type Tcp struct {
	wg              *sync.WaitGroup
	mutex           *sync.Mutex
	ip              net.IPAddr
	tcpPort         int
	pingCount       int
	csv             *[]utils.CloudflareIPData
	control         chan bool
	progressHandler func(e utils.ProgressEvent)
}
