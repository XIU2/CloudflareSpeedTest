package utils

import "github.com/cheggaaa/pb/v3"

type ProgressEvent int

const (
	NoAvailableIPFound ProgressEvent = iota
	AvailableIPFound
	NormalPing
)

type Bar struct {
	*pb.ProgressBar
}

func NewBar(count int) *Bar {
	return &Bar{pb.Simple.Start(count)}
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
