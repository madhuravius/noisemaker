package internal

import (
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	"github.com/urfave/cli/v2"
)

type RunConfig struct {
	cpuFuncs         []func(wg *sync.WaitGroup)
	data             []int64
	socketData       []byte
	desiredCpu       float64
	desiredMem       float64
	desiredPort      int
	desiredBandwidth int
}

func NewRunConfig(desiredCpu, desiredMem float64, desiredPort, desiredBandwidth int) *RunConfig {
	r := &RunConfig{
		data:             []int64{},
		desiredCpu:       desiredCpu,
		desiredMem:       desiredMem,
		desiredPort:      desiredPort,
		desiredBandwidth: desiredBandwidth,
	}
	r.generateSocketDataToSend()
	go r.startCpuUsingFuncs()
	go r.startHttpServer()
	return r
}

func (r *RunConfig) printDebugStatement() {
	usageMessage := fmt.Sprintf("Running: Memory - %f, CPU - %f (%d goroutines).",
		r.getUsedMemory(),
		r.getUsedCpu(),
		len(r.cpuFuncs),
	)

	// begin - if not really blowing up machine for other reasons, print it here
	if len(r.data) == 0 || len(r.cpuFuncs) == 0 {
		usageMessage += " WARNINGS: "
	}
	if len(r.data) == 0 {
		usageMessage += "not allocating memory"
	}
	if len(r.cpuFuncs) == 0 {
		if len(r.data) == 0 {
			usageMessage += ","
		}
		usageMessage += "not allocating cpu"
	}
	if len(r.data) == 0 || len(r.cpuFuncs) == 0 {
		usageMessage += " because using more outside application"
	}
	// end - if not really blowing up machine for other reasons, print it here

	log.Printf("%s\n", usageMessage)
}
func RunLoop(ctx *cli.Context) {
	stop := make(chan bool)

	desiredCpu := math.Min(ctx.Value(CpuContextKey).(float64), 99.)
	desiredMem := math.Min(ctx.Value(MemoryContextKey).(float64), 99.)
	desiredPort := ctx.Value(PortContextKey).(int)
	desiredBandwidth := ctx.Value(BandwidthContextKey).(int)

	log.Printf("Running with requested values: memory - %f, cpu - %f, port - %d (bandwidth: %d)", desiredMem, desiredCpu, desiredPort, desiredBandwidth)
	r := NewRunConfig(desiredCpu, desiredMem, desiredPort, desiredBandwidth)
	go func(r *RunConfig) {
		for {
			r.calculateAndUseMem()
			r.calculateAndUseCPU()
			r.trashbinClientSocketFunc()
		}
	}(r)
	for {
		select {
		case <-time.After(DebugLogInterval):
			r.printDebugStatement()
		case <-stop:
			return
		}
	}
}
