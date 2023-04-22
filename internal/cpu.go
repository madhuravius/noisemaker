package internal

import (
	"sync"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

const (
	FibSize = 10
)

func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

// loadCPUFunc - run CPU load in specify cores count and percentage
func loadCPUFunc(wg *sync.WaitGroup) {
	defer wg.Done()
	timeStarted := time.Now()
	for {
		if time.Now().After(timeStarted.Add(CPUPollingTime)) {
			break
		}
		_ = fibonacci(FibSize)
		time.Sleep(10 * time.Microsecond)
	}
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
			for _, fn := range r.cpuFuncs {
				wg.Add(1)
				go fn(&wg)
			}
			wg.Wait()
		}
	}
}
