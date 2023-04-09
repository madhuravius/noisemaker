package internal

import (
	"sync"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

// loadCPUFunc - run CPU load in specify cores count and percentage
// adapted somewhat from: https://stackoverflow.com/a/41084841
// TODO - may want to write a better more predictable CPU user since this is very all/nothing
func loadCPUFunc(wg *sync.WaitGroup) {
	defer wg.Done()
	done := make(chan int)
	go func() {
		for {
			select {
			case <-done:
				return
			//nolint:staticcheck
			default:
			}
		}
	}()
	time.Sleep(CPUPollingTime)
	close(done)
}

func (r *RunConfig) getUsedCpu() float64 {
	cpuStats, err := cpu.Percent(CPUPollingTime, false)
	if err != nil {
		panic(err)
	}
	cpuUsage := 0.
	for _, cpuThread := range cpuStats {
		cpuUsage += cpuThread
	}
	return cpuUsage / float64(len(cpuStats))
}
func (r *RunConfig) calculateAndUseCPU() {
	usedCpu := r.getUsedCpu()
	if usedCpu <= r.desiredCpu {
		r.cpuFuncs = append(r.cpuFuncs, loadCPUFunc)
	} else if len(r.cpuFuncs) > 0 {
		r.cpuFuncs = r.cpuFuncs[:len(r.cpuFuncs)-1]
	}
}

func (r *RunConfig) startCpuUsingFuncs() {
	for {
		if len(r.cpuFuncs) > 0 {
			wg := sync.WaitGroup{}
			for _, cpuFunc := range r.cpuFuncs {
				wg.Add(1)
				go cpuFunc(&wg)
			}
			wg.Wait()
		}
	}
}
