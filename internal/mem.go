package internal

import "github.com/shirou/gopsutil/mem"

func (r *RunConfig) getUsedMemory() float64 {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		panic(err)
	}
	return vmStat.UsedPercent
}

func (r *RunConfig) calculateAndUseMem() {
	if r.getUsedMemory() <= r.desiredMem {
		// arbitrary number chosen, scaled up to something that seemed to move fast enough, and drains
		for i := 0; i < ArrayChunkSize*ArrayChunkSize*25; i++ {
			r.data = append(r.data, 1)
		}
	} else if r.getUsedMemory() > r.desiredMem && len(r.data) > ArrayChunkSize {
		_, r.data = r.data[0], r.data[ArrayChunkSize:]
	}
}
