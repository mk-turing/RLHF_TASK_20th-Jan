package monitor

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"time"
)

type ResourceStats struct {
	CPUUsage    float64
	MemoryUsage uint64
}

func MonitorResources(interval time.Duration, stopChan <-chan struct{}, statsChan chan<- ResourceStats) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-stopChan:
			return
		case <-ticker.C:
			cpuPercentages, _ := cpu.Percent(0, false)
			virtualMem, _ := mem.VirtualMemory()
			statsChan <- ResourceStats{
				CPUUsage:    cpuPercentages[0],
				MemoryUsage: virtualMem.Used,
			}
		}
	}
}

func PrintResourceUsage(stats ResourceStats) {
	fmt.Printf("CPU Usage: %.2f%%, Memory Usage: %dMB\n", stats.CPUUsage, stats.MemoryUsage/(1024*1024))
}
