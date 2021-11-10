package utils

import "github.com/cheggaaa/pb/v3"

type ProgressEvent int

const (
	NoAvailableIPFound ProgressEvent = iota
	AvailableIPFound
	NormalPing
)

type Bar struct {
	pb *pb.ProgressBar
}

func NewBar(count int) *Bar {
	return &Bar{pb.Simple.Start(count)}
}

func (b *Bar) Grow(num int) {
	b.pb.Add(num)
}

func (b *Bar) Done() {
	b.pb.Finish()
}

func handleProgressGenerator(pb *pb.ProgressBar) func(e ProgressEvent) {
	return func(e ProgressEvent) {
		switch e {
		case NoAvailableIPFound:
			// pb.Add(pingTime)
		case AvailableIPFound:
			// pb.Add(failTime)
		case NormalPing:
			pb.Increment()
		}
	}
}
