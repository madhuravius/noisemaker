package internal

import "time"

const (
	ArrayChunkSize      = 1024
	CpuContextKey       = "cpu"
	MemoryContextKey    = "mem"
	BandwidthContextKey = "bandwidth"
	PortContextKey      = "port"
	DebugLogInterval    = time.Second * 3
	CPUPollingTime      = time.Millisecond * 100
	FileSizeToSend      = 100 // MB
)
