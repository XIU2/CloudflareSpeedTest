package utils

type ProgressEvent int

const (
	NoAvailableIPFound ProgressEvent = iota
	AvailableIPFound
	NormalPing
)
