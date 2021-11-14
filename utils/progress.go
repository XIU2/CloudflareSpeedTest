package utils

import "github.com/cheggaaa/pb/v3"

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